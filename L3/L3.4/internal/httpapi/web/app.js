const uploadForm = document.getElementById("uploadForm");
const imagesList = document.getElementById("imagesList");
const refreshBtn = document.getElementById("refreshBtn");

const statusLabels = {
  pending: "В очереди",
  processing: "Обрабатывается",
  ready: "Готово",
  failed: "Ошибка",
};

uploadForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const formData = new FormData(uploadForm);
  try {
    const resp = await fetch("/upload", {
      method: "POST",
      body: formData,
    });
    if (!resp.ok) {
      const payload = await resp.json().catch(() => ({}));
      throw new Error(payload.error || "Не удалось загрузить файл");
    }
    uploadForm.reset();
    await loadImages();
  } catch (err) {
    alert(err.message);
  }
});

refreshBtn.addEventListener("click", () => loadImages());

async function loadImages() {
  imagesList.innerHTML = "<p>Загружаем...</p>";
  try {
    const resp = await fetch("/images?limit=100");
    const items = await resp.json();
    renderImages(items);
  } catch (err) {
    imagesList.innerHTML = `<p>Ошибка загрузки: ${err.message}</p>`;
  }
}

function renderImages(items) {
  if (!Array.isArray(items) || items.length === 0) {
    imagesList.innerHTML = "<p>Пока ничего нет — загрузите первую фотографию.</p>";
    return;
  }
  imagesList.innerHTML = "";
  items.forEach((item) => {
    const card = document.createElement("article");
    card.className = "image-card";

    const preview = document.createElement("img");
    preview.className = "image-preview";
    preview.alt = item.image_id;
    if (item.processed_url) {
      preview.src = item.processed_url;
    } else if (item.original_url) {
      preview.src = item.original_url;
    } else {
      preview.src = "data:image/gif;base64,R0lGODlhAQABAAAAACw=";
    }
    card.appendChild(preview);

    const statusBadge = document.createElement("span");
    statusBadge.className = `badge ${item.status}`;
    statusBadge.textContent = statusLabels[item.status] || item.status;
    card.appendChild(statusBadge);

    if (item.error_message) {
      const err = document.createElement("p");
      err.className = "hint";
      err.style.color = "#f87171";
      err.textContent = item.error_message;
      card.appendChild(err);
    }

    const actions = document.createElement("div");
    actions.className = "card-actions";

    const viewOriginal = document.createElement("a");
    viewOriginal.textContent = "Оригинал";
    viewOriginal.href = `/image/${item.image_id}?variant=original`;
    viewOriginal.target = "_blank";
    actions.appendChild(viewOriginal);

    const viewProcessed = document.createElement("a");
    viewProcessed.textContent = "Результат";
    viewProcessed.href = `/image/${item.image_id}?variant=processed`;
    viewProcessed.target = "_blank";
    viewProcessed.disabled = !item.processed_url;
    if (!item.processed_url) {
      viewProcessed.classList.add("disabled");
      viewProcessed.href = "#";
    }
    actions.appendChild(viewProcessed);

    const removeBtn = document.createElement("button");
    removeBtn.type = "button";
    removeBtn.textContent = "Удалить";
    removeBtn.addEventListener("click", () => deleteImage(item.image_id));
    actions.appendChild(removeBtn);

    card.appendChild(actions);
    imagesList.appendChild(card);
  });
}

async function deleteImage(imageId) {
  if (!confirm("Удалить изображение?")) {
    return;
  }
  await fetch(`/image/${imageId}`, { method: "DELETE" });
  await loadImages();
}

loadImages().catch(() => {
  imagesList.innerHTML = "<p>Не удалось загрузить список изображений.</p>";
});

