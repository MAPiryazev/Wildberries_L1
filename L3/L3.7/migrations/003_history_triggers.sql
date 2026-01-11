-- Migration: Add triggers for item history tracking

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION log_item_created() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO item_history (item_id, changed_by, action, changes)
    VALUES (
        NEW.id,
        NEW.created_by,
        'created',
        jsonb_build_object(
            'name', NEW.name,
            'sku', NEW.sku,
            'quantity', NEW.quantity,
            'location', NEW.location
        )
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION log_item_update()
RETURNS TRIGGER AS $$
DECLARE
    changes JSONB := '{}'::jsonb;
BEGIN
    IF OLD.name IS DISTINCT FROM NEW.name THEN
        changes := jsonb_set(changes, '{name}', jsonb_build_object('old', OLD.name, 'new', NEW.name));
    END IF;

    IF OLD.sku IS DISTINCT FROM NEW.sku THEN
        changes := jsonb_set(changes, '{sku}', jsonb_build_object('old', OLD.sku, 'new', NEW.sku));
    END IF;

    IF OLD.quantity IS DISTINCT FROM NEW.quantity THEN
        changes := jsonb_set(changes, '{quantity}', jsonb_build_object('old', OLD.quantity, 'new', NEW.quantity));
    END IF;

    IF OLD.location IS DISTINCT FROM NEW.location THEN
        changes := jsonb_set(changes, '{location}', jsonb_build_object('old', OLD.location, 'new', NEW.location));
    END IF;

    IF changes != '{}'::jsonb THEN
        INSERT INTO item_history (item_id, changed_by, action, changes)
        VALUES (NEW.id, NEW.updated_by, 'updated', changes);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION log_item_deleted()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO item_history (item_id, changed_by, action, changes)
    VALUES (
        OLD.id,
        OLD.updated_by,
        'deleted',
        jsonb_build_object(
            'name', OLD.name,
            'sku', OLD.sku,
            'quantity', OLD.quantity,
            'location', OLD.location
        )
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

DROP TRIGGER IF EXISTS trigger_item_created ON items;
DROP TRIGGER IF EXISTS trigger_item_updated ON items;
DROP TRIGGER IF EXISTS trigger_item_deleted ON items;

CREATE TRIGGER trigger_item_created
    AFTER INSERT ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_created();

CREATE TRIGGER trigger_item_updated
    AFTER UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_update();

CREATE TRIGGER trigger_item_deleted
    BEFORE DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_deleted();
