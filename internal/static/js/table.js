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

  document.getElementById("add-menu").classList.toggle("active", isOpenAddMenu);

  document.getElementById("cart").classList.toggle("active", isOpenCart);

  document.body.classList.toggle(
    "overflow-hidden",
    isOpenAddMenu || isOpenCart,
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
