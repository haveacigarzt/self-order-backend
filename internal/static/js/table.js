document.addEventListener("DOMContentLoaded", () => {
  const elements = document.querySelectorAll("[data-harga]");
  const formatter = new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  });
  elements.forEach((element) => {
    const harga = Number(element.dataset.harga);
    const formatted = formatter.format(harga);
    element.textContent = formatted;
  });
});

const navContainer = document.getElementById("nav-container");
const navLinks = navContainer.querySelectorAll('a[href^="#grup-"]');
const sections = document.querySelectorAll('[id^="grup-"]');

let active = "";

const observer = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        active = entry.target.id;

        const activeLink = navContainer.querySelector(`[href="#${active}"]`);

        activeLink?.scrollIntoView({
          behavior: "smooth",
          block: "nearest",
          inline: "center",
        });

        navLinks.forEach((link) => {
          link.classList.toggle(
            "active",
            link.getAttribute("href") === `#${active}`,
          );
        });
      }
    });
  },
  {
    rootMargin: "-64px 0px -60% 0px",
    threshold: 0,
  },
);

sections.forEach((section) => {
  observer.observe(section);
});

const buttonBackToTop = document.getElementById("button_back_to_top");
buttonBackToTop.addEventListener("click", () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
});

const backToTop = document.getElementById("back_to_top");

const checkScroll = () => {
  const show =
    window.innerHeight + window.scrollY >=
    document.documentElement.scrollHeight - 50;

  backToTop.classList.toggle("hidden", !show);
};

window.addEventListener("scroll", checkScroll);
window.addEventListener("resize", checkScroll);

checkScroll();

// Modal detail produk
const produkCards = document.querySelectorAll("[data-produk]");
produkCards.forEach((card) => {
  card.addEventListener("click", () => {
    const encoded = card.dataset.produk;

    const json = atob(encoded);
    const produk = JSON.parse(json);

    // console.log(produk);
  });
});

const modal = document.getElementById("modal_produk");
const form = document.getElementById("add_menu_form");

const closeButton = document.getElementById("button_close_modal");

const produkNama = document.getElementById("modal_produk_nama");
const produkDeskripsi = document.getElementById("modal_produk_deskripsi");
const produkHarga = document.getElementById("modal_produk_harga");
const quantityElement = document.getElementById("modal_produk_quantity");
const produkThumbnail = document.getElementById("modal_produk_image");
const produkIdInput = document.getElementById("produk_id_input");
const quantityInput = document.getElementById("quantity_input");
const notesInput = document.getElementById("notes_input");
const hapusButton = document.getElementById("hapus_button");
const simpanButton = document.getElementById("simpan_button");

const totalQuantityElement = document.getElementById("total_quantity");
const totalAmountElement = document.getElementById("total_amount");

let selectedProduk = null;

// Format Rupiah
const formatRupiah = (harga) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(harga);
};

let quantity = 0;

form.addEventListener("submit", () => {
  simpanButton.disabled = true;
});

// Increment
document.getElementById("button_increment").addEventListener("click", () => {
  quantity++;

  quantityElement.textContent = quantity;
  quantityInput.value = quantity;

  simpanButton.disabled = quantity === selectedProduk.qty;
});

// Decrement
document.getElementById("button_decrement").addEventListener("click", () => {
  if (quantity <= 0) return;

  quantity--;
  quantityElement.textContent = quantity;
  quantityInput.value = quantity;

  simpanButton.disabled = quantity === selectedProduk.qty || quantity === 0;
});

notesInput.addEventListener("input", (e) => {
  if (selectedProduk.qty === 0) return;
  simpanButton.disabled = e.target.value === selectedProduk.catatan;
});

hapusButton.addEventListener("click", async () => {
  hapusButton.disabled = true;
  const noMeja = hapusButton.dataset.meja;
  try {
    const response = await fetch(`/meja/${noMeja}/${selectedProduk.id}`, {
      method: "DELETE",
    });
    if (!response.ok) {
      throw new Error("Gagal menghapus pesanan");
    }
    window.location.href = response.url;
  } catch (error) {
    console.error(error);
    alert("Gagal menghapus pesanan. Silakan coba kembali.");
  } finally {
    hapusButton.disabled = false;
  }
});

//
const addMenu = document.getElementById("add-menu");

function openAddMenu(produk) {
  console.log("open addmenu");
  const url = new URL(window.location.href);

  url.searchParams.set("id", produk.id);
  url.hash = "add-menu";

  history.pushState({}, "", url);

  addMenu.classList.add("active");
  document.body.classList.add("overflow-hidden");

  selectedProduk = produk;
  quantity = selectedProduk.qty;

  // Isi informasi produk
  produkNama.textContent = selectedProduk.nama;
  produkDeskripsi.textContent = selectedProduk.deskripsi;
  produkHarga.textContent = formatRupiah(selectedProduk.harga);
  produkIdInput.value = selectedProduk.id;
  if (selectedProduk.gambar[0].filepath) {
    produkThumbnail.src = `/static/uploads/${selectedProduk.gambar[0].filepath}`;
    produkThumbnail.alt = "product image";
  } else {
    produkThumbnail.src = "/static/img/placeholder.svg";
    produkThumbnail.alt = "placeholder image";
  }

  quantityElement.textContent = selectedProduk.qty;
  quantityInput.value = selectedProduk.qty;
  notesInput.value = selectedProduk.catatan;

  if (selectedProduk.qty == 0) {
    hapusButton.disabled = true;
  }
}

function closeAddMenu() {
  console.log("closing addmenu");

  document.getElementById("add-menu").classList.remove("active");

  document.body.classList.remove("overflow-hidden");

  history.replaceState(
    {},
    "",
    window.location.pathname + window.location.search,
  );
}

function closeCart() {
  console.log("closing cart");

  document.getElementById("cart").classList.remove("active");
  document.body.classList.remove("overflow-hidden");

  history.replaceState(
    {},
    "",
    window.location.pathname + window.location.search,
  );
}

document.getElementById("close_page").addEventListener("click", closeAddMenu);

window.addEventListener("popstate", () => {
  const url = new URL(window.location.href);

  const isOpenAddMenu = url.hash === "#add-menu";
  const isOpenCart = url.hash === "#cart";
  const isOpenCarousel = url.hash === "#carousel";

  document.getElementById("add-menu").classList.toggle("active", isOpenAddMenu);

  document.getElementById("cart").classList.toggle("active", isOpenCart);

  document
    .getElementById("modal_image_carousel")
    .classList.toggle("active", isOpenCarousel);

  document.body.classList.toggle(
    "overflow-hidden",
    isOpenAddMenu || isOpenCart || isOpenCarousel,
  );
});

document.querySelectorAll(".produk-card").forEach((card) => {
  card.addEventListener("click", () => {
    try {
      const encoded = card.dataset.produk;
      const produk = JSON.parse(atob(encoded));

      openAddMenu(produk);
    } catch (error) {
      console.error("Gagal membaca data produk:", error);
    }
  });
});

//
const cart = document.getElementById("cart");

function openCart() {
  console.log("open");
  const url = new URL(window.location.href);

  url.hash = "cart";

  history.pushState({}, "", url);

  cart.classList.add("active");
  document.body.classList.add("overflow-hidden");
}

document.getElementById("close_cart").addEventListener("click", closeCart);
document.getElementById("add_item_cart").addEventListener("click", closeCart);
document.getElementById("cart_button").addEventListener("click", openCart);

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

document.querySelectorAll(".edit-button").forEach((button) => {
  button.addEventListener("click", () => {
    try {
      closeCart();
      const encoded = button.dataset.produk;
      const produk = JSON.parse(atob(encoded));
      console.log("Start");
      sleep(300).then(() => {
        openAddMenu(produk);
        console.log(produk);
        console.log("Executed after 1.5 seconds");
      });
    } catch (error) {
      console.error("Gagal membaca data produk:", error);
    }
  });
});

//
const modal_image_carousel = document.getElementById("modal_image_carousel");

const markers = document.getElementById("carousel-markers");

function openImageCarousel() {
  console.log("open image carousel");
  const url = new URL(window.location.href);

  url.hash = "carousel";

  history.pushState({}, "", url);

  modal_image_carousel.classList.add("active");
  modal_image_carousel.parentElement.classList.add("overflow-hidden");

  const carousel = document.querySelector(".carousel");
  const track = document.getElementById("carousel-track");

  const prev = document.getElementById("prev");
  const next = document.getElementById("next");

  let inner = "";

  selectedProduk.gambar.forEach((gambar) => {
    if (!gambar.filepath) return;

    inner += `
    <div class="card">
      <img
        src="/static/uploads/${gambar.filepath}"
        alt=""
        class="aspect-square object-cover"
        draggable="false"
      >
    </div>`;
  });

  track.innerHTML = inner;

  const cards = track.querySelectorAll(".card");

  let currentIndex = 0;

  let startX = 0;
  let currentX = 0;
  let isDragging = false;
  let activePointerId = null;

  // =================================
  // 1. Buat marker
  // =================================

  cards.forEach((_, index) => {
    const marker = document.createElement("button");

    marker.type = "button";
    marker.className = "carousel-marker";

    marker.addEventListener("click", () => {
      goToSlide(index);
    });

    markers.appendChild(marker);
  });

  const markerItems = markers.querySelectorAll(".carousel-marker");

  // =================================
  // 2. Navigasi slide
  // =================================

  function goToSlide(index) {
    currentIndex = Math.max(0, Math.min(index, cards.length - 1));

    track.style.transform = `translateX(-${currentIndex * 100}%)`;

    markerItems.forEach((marker, i) => {
      marker.classList.toggle("active", i === currentIndex);
    });

    prev.disabled = currentIndex === 0;
    next.disabled = currentIndex === cards.length - 1;
  }

  // =================================
  // 3. Tombol navigasi
  // =================================

  prev.addEventListener("click", () => {
    goToSlide(currentIndex - 1);
  });

  next.addEventListener("click", () => {
    goToSlide(currentIndex + 1);
  });

  // =================================
  // 4. Mulai drag
  // =================================

  carousel.addEventListener("pointerdown", (e) => {
    if (e.pointerType === "mouse" && e.button !== 0) return;
    if (isDragging) return;

    startX = e.clientX;
    currentX = e.clientX;

    isDragging = true;
    activePointerId = e.pointerId;

    carousel.setPointerCapture(e.pointerId);

    track.classList.add("dragging");
  });

  // =================================
  // 5. Gerakkan carousel
  // =================================

  carousel.addEventListener("pointermove", (e) => {
    if (!isDragging || e.pointerId !== activePointerId) return;

    currentX = e.clientX;

    const distance = currentX - startX;
    const width = carousel.offsetWidth;

    let offset = -(currentIndex * width) + distance;

    // Batasi pergerakan pada awal dan akhir carousel
    if (currentIndex === 0 && distance > 0) {
      offset = distance * 0.3;
    }

    if (currentIndex === cards.length - 1 && distance < 0) {
      offset = -(currentIndex * width) + distance * 0.3;
    }

    track.style.transform = `translateX(${offset}px)`;
  });

  // =================================
  // 6. Selesaikan drag
  // =================================

  function endDrag(e) {
    if (!isDragging || e.pointerId !== activePointerId) return;

    const distance = currentX - startX;
    const threshold = carousel.offsetWidth * 0.2;

    isDragging = false;
    activePointerId = null;

    track.classList.remove("dragging");

    if (distance < -threshold) {
      goToSlide(currentIndex + 1);
    } else if (distance > threshold) {
      goToSlide(currentIndex - 1);
    } else {
      goToSlide(currentIndex);
    }
  }

  carousel.addEventListener("pointerup", endDrag);
  carousel.addEventListener("pointercancel", endDrag);

  // =================================
  // 7. Inisialisasi
  // =================================

  if (cards.length > 0) {
    goToSlide(0);
  } else {
    prev.disabled = true;
    next.disabled = true;
  }
}

function closeImageCarousel() {
  console.log("close image carousel");

  modal_image_carousel.classList.remove("active");
  modal_image_carousel.parentElement.classList.remove("overflow-hidden");

  history.replaceState(
    {},
    "",
    window.location.pathname + window.location.search,
  );
  markers.replaceChildren();
}

produkThumbnail.addEventListener("click", (e) => {
  const target = e.target;
  if (target.alt == "placeholder image") return;
  // console.log(selectedProduk);
  openImageCarousel();
});

const button_max_image = document.getElementById("button_max_image");
button_max_image.addEventListener("click", (e) => {
  const target = e.currentTarget.previousElementSibling;
  if (target.alt == "placeholder image") return;
  // console.log(selectedProduk);
  openImageCarousel();
});

document
  .getElementById("close_image_carousel")
  .addEventListener("click", closeImageCarousel);

const pembayaran = document.getElementById("pembayaran");

document
  .getElementById("pembayaran_button")
  .addEventListener("click", openPembayaran);

document
  .getElementById("close_pembayaran")
  .addEventListener("click", closePembayaran);

function openPembayaran() {
  console.log("open pembayaran");
  const url = new URL(window.location.href);

  url.hash = "pembayaran";

  history.pushState({}, "", url);

  pembayaran.classList.add("active");
  pembayaran.parentElement.classList.add("overflow-hidden");
}

function closePembayaran() {
  console.log("close pembayaran");

  pembayaran.classList.remove("active");
  pembayaran.parentElement.classList.remove("overflow-hidden");

  history.replaceState(
    {},
    "",
    window.location.pathname + window.location.search,
  );
}
