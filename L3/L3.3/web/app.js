const API_BASE = '/api';

async function fetchJSON(url, opts = {}) {
  const res = await fetch(url, { headers: { 'Content-Type': 'application/json' }, ...opts });
  if (!res.ok) {
    let msg = `${res.status} ${res.statusText}`;
    try { const e = await res.json(); if (e && e.error) msg = e.error; } catch {}
    throw new Error(msg);
  }
  if (res.status === 204) return null;
  return res.json();
}

function renderTree(container, nodes, depth = 0) {
  container.innerHTML = '';
  const ul = document.createElement('ul');
  nodes.forEach(node => {
    const li = document.createElement('li');
    li.style.marginLeft = `${depth * 16}px`;
    const header = document.createElement('div');
    header.className = 'comment';
    header.innerHTML = `#${node.id} <strong>user ${node.user_id}</strong> · <span class="date">${new Date(node.created_at).toLocaleString()}</span><br>${escapeHtml(node.content)}`;

    const actions = document.createElement('div');
    actions.className = 'actions';

    const delBtn = document.createElement('button');
    delBtn.textContent = 'Удалить';
    delBtn.addEventListener('click', async () => {
      const userId = prompt('Введите ваш user_id для удаления:');
      if (!userId) return;
      try {
        await fetchJSON(`${API_BASE}/comments/${node.id}?user_id=${encodeURIComponent(userId)}`, { method: 'DELETE' });
        await loadRoot();
      } catch (e) { alert(e.message); }
    });

    const replyBtn = document.createElement('button');
    replyBtn.textContent = 'Ответить';
    replyBtn.addEventListener('click', () => {
      document.querySelector('#parentId').value = node.id;
      document.querySelector('#content').focus();
    });

    actions.appendChild(replyBtn);
    actions.appendChild(delBtn);
    header.appendChild(actions);

    li.appendChild(header);

    if (node.children && node.children.length) {
      const childContainer = document.createElement('div');
      renderTree(childContainer, node.children, depth + 1);
      li.appendChild(childContainer);
    }

    ul.appendChild(li);
  });
  container.appendChild(ul);
}

function escapeHtml(s) {
  return String(s)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

async function loadRoot() {
  const root = document.querySelector('#commentsRoot');
  try {
    const trees = await fetchJSON(`${API_BASE}/comments`);
    renderTree(root, trees);
  } catch (e) {
    root.innerHTML = `<div class="error">${escapeHtml(e.message)}</div>`;
  }
}

async function createComment() {
  const userId = parseInt(document.querySelector('#userId').value, 10);
  const parentIdStr = document.querySelector('#parentId').value.trim();
  const content = document.querySelector('#content').value.trim();
  const msg = document.querySelector('#createMsg');
  msg.textContent = '';
  try {
    const body = { user_id: userId, content };
    if (parentIdStr) body.parent_id = parseInt(parentIdStr, 10);
    await fetchJSON(`${API_BASE}/comments`, { method: 'POST', body: JSON.stringify(body) });
    document.querySelector('#content').value = '';
    await loadRoot();
    msg.textContent = 'Создано';
    setTimeout(() => { msg.textContent = ''; }, 2000);
  } catch (e) {
    msg.textContent = e.message;
  }
}

async function search() {
  const q = document.querySelector('#searchInput').value.trim();
  const ul = document.querySelector('#searchResults');
  ul.innerHTML = '';
  if (q.length < 3) {
    ul.innerHTML = '<li>Минимум 3 символа</li>';
    return;
  }
  try {
    const items = await fetchJSON(`${API_BASE}/comments/search?query=${encodeURIComponent(q)}&limit=20`);
    if (!items.length) {
      ul.innerHTML = '<li>Ничего не найдено</li>';
      return;
    }
    for (const it of items) {
      const li = document.createElement('li');
      li.innerHTML = `#${it.id} · user ${it.user_id} · ${new Date(it.created_at).toLocaleString()} — ${escapeHtml(it.content)}`;
      ul.appendChild(li);
    }
  } catch (e) {
    ul.innerHTML = `<li class="error">${escapeHtml(e.message)}</li>`;
  }
}

document.addEventListener('DOMContentLoaded', () => {
  document.querySelector('#createBtn').addEventListener('click', createComment);
  document.querySelector('#searchBtn').addEventListener('click', search);
  loadRoot();
});


