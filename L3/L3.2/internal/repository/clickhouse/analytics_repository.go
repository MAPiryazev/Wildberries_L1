package clickhouse

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"shortener/internal/models"
)

// HTTPClient абстракция для тестов
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AnalyticsRepo — ClickHouse через HTTP API
type AnalyticsRepo struct {
	baseURL  string // http://clickhouse:8123
	database string // shortener_analytics
	user     string
	password string
	client   HTTPClient
}

func NewAnalyticsRepo(host string, httpPort string, database string, user string, password string, client HTTPClient) *AnalyticsRepo {
	base := fmt.Sprintf("http://%s:%s", host, httpPort)
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &AnalyticsRepo{baseURL: base, database: database, user: user, password: password, client: client}
}

func (r *AnalyticsRepo) InsertClick(ctx context.Context, e models.ClickEvent) error {
	q := "INSERT INTO " + r.database + ".click_analytics (short_code, client_id, user_agent, ip, timestamp) VALUES"
	esc := func(s string) string { return strings.ReplaceAll(s, "'", "\\'") }
	payload := fmt.Sprintf("('%s','%s','%s','%s',toDateTime(%d))",
		esc(e.ShortCode), e.ClientID.String(), esc(e.UserAgent), esc(e.IP), e.At.Unix())
	return r.exec(ctx, q+" "+payload)
}

func (r *AnalyticsRepo) CountByShortCode(ctx context.Context, shortCode string) (int64, error) {
	q := fmt.Sprintf("SELECT count() FROM %s.click_analytics WHERE short_code = '%s'", r.database, escape(shortCode))
	return r.querySingleInt64(ctx, q)
}

func (r *AnalyticsRepo) AggregateDaily(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	q := fmt.Sprintf(`
		SELECT formatDateTime(timestamp, '%%Y-%%m-%%d') AS day, count() AS c
		FROM %s.click_analytics
		WHERE short_code = '%s'
		GROUP BY day
		ORDER BY day ASC
	`, r.database, escape(shortCode))
	return r.queryAgg(ctx, q)
}

func (r *AnalyticsRepo) AggregateByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	q := fmt.Sprintf(`
		SELECT user_agent AS ua, count() AS c
		FROM %s.click_analytics
		WHERE short_code = '%s'
		GROUP BY ua
		ORDER BY c DESC
	`, r.database, escape(shortCode))
	return r.queryAgg(ctx, q)
}

// ---- helpers ----

func (r *AnalyticsRepo) exec(ctx context.Context, sql string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL, bytes.NewBufferString(sql))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.SetBasicAuth(r.user, r.password)
	q := req.URL.Query()
	q.Set("database", r.database)
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("clickhouse exec %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (r *AnalyticsRepo) querySingleInt64(ctx context.Context, sql string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL, bytes.NewBufferString(sql))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.SetBasicAuth(r.user, r.password)
	q := req.URL.Query()
	q.Set("database", r.database)
	q.Set("default_format", "CSV")
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("clickhouse query %d: %s", resp.StatusCode, string(b))
	}
	rdr := csv.NewReader(resp.Body)
	rec, err := rdr.Read()
	if err != nil {
		return 0, err
	}
	if len(rec) == 0 {
		return 0, fmt.Errorf("empty result")
	}
	var v int64
	_, err = fmt.Sscan(rec[0], &v)
	return v, err
}

func (r *AnalyticsRepo) queryAgg(ctx context.Context, sql string) ([]models.AggPoint, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL, bytes.NewBufferString(sql))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.SetBasicAuth(r.user, r.password)
	q := req.URL.Query()
	q.Set("database", r.database)
	q.Set("default_format", "CSV")
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("clickhouse query %d: %s", resp.StatusCode, string(b))
	}
	rdr := csv.NewReader(resp.Body)
	var out []models.AggPoint
	for {
		rec, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 2 {
			continue
		}
		var c int64
		if _, err := fmt.Sscan(rec[1], &c); err != nil {
			return nil, err
		}
		out = append(out, models.AggPoint{Key: rec[0], Count: c})
	}
	return out, nil
}

func escape(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}
