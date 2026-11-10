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
const modalContent = document.getElementById("modal_content");

const closeButton = document.getElementById("button_close_modal");

const produkNama = document.getElementById("modal_produk_nama");
const produkDeskripsi = document.getElementById("modal_produk_deskripsi");
const produkHarga = document.getElementById("modal_produk_harga");

const quantityElement = document.getElementById("modal_produk_quantity");
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

const pesanan = [];

let totalQuantity = 0;
let totalAmount = 0;

// Open Modal
function openModal(produk) {
  selectedProduk = produk;

  // Isi informasi produk
  produkNama.textContent = produk.nama;
  produkDeskripsi.textContent = produk.deskripsi;
  produkHarga.textContent = formatRupiah(produk.harga);

  const p = pesanan.find((e) => e.id === selectedProduk.id);
  if (p) {
    quantityElement.textContent = p.quantity;
    notesInput.value = p.notes;
  } else {
    quantityElement.textContent = 0;
    notesInput.value = "";
  }

  modal.classList.remove("hidden");
  document.body.classList.add("overflow-hidden");
}

function closeModal() {
  modal.classList.add("hidden");
  document.body.classList.remove("overflow-hidden");
  selectedProduk = null;
}

// Product Card Click
document.querySelectorAll(".produk-card").forEach((card) => {
  card.addEventListener("click", () => {
    const encoded = card.dataset.produk;
    const produk = JSON.parse(atob(encoded));

    openModal(produk);
  });
});

// Close Button
closeButton.addEventListener("click", closeModal);
simpanButton.addEventListener("click", closeModal);

// Click Outside
modal.addEventListener("click", (event) => {
  if (event.target === modal) {
    closeModal();
  }
});

// Escape
document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && !modal.classList.contains("hidden")) {
    closeModal();
  }
});

// Increment
document.getElementById("button_increment").addEventListener("click", () => {
  let p = pesanan.find((e) => e.id === selectedProduk.id);

  if (p) {
    p.quantity++;
  } else {
    p = {
      ...selectedProduk,
      quantity: 1,
      notes: document.getElementById("notes_input").value,
    };
    pesanan.push(p);
  }

  quantityElement.textContent = p.quantity;

  totalQuantity++;
  if (totalQuantityElement.classList.contains("hidden")) {
    totalQuantityElement.classList.replace("hidden", "flex");
  }
  totalQuantityElement.textContent = totalQuantity;

  totalAmount = totalAmount + selectedProduk.harga;
  totalAmountElement.textContent = formatRupiah(totalAmount);

  console.log(pesanan);
});

// Decrement
document.getElementById("button_decrement").addEventListener("click", () => {
  const p = pesanan.find((e) => e.id === selectedProduk.id);

  if (!p) return;

  p.quantity--;

  if (p.quantity <= 0) {
    const index = pesanan.indexOf(p);
    pesanan.splice(index, 1);

    quantityElement.textContent = 0;
  } else {
    quantityElement.textContent = p.quantity;
  }

  totalQuantity--;
  if (totalQuantity == 0) {
    totalQuantityElement.classList.replace("flex", "hidden");
  }
  totalQuantityElement.textContent = totalQuantity;

  totalAmount = totalAmount - selectedProduk.harga;
  totalAmountElement.textContent = formatRupiah(totalAmount);

  console.log(pesanan);
});

notesInput.addEventListener("input", () => {
  const p = pesanan.find((e) => e.id === selectedProduk.id);
  if (!p) return;
  p.notes = event.target.value;
});

hapusButton.addEventListener("click", () => {
  const p = pesanan.find((e) => e.id === selectedProduk.id);
  if (!p) return;
  totalQuantity = totalQuantity - p.quantity;
  if (totalQuantity == 0) {
    totalQuantityElement.classList.replace("flex", "hidden");
  }
  totalQuantityElement.textContent = totalQuantity;

  totalAmount = totalAmount - selectedProduk.harga * p.quantity;
  totalAmountElement.textContent = formatRupiah(totalAmount);

  const index = pesanan.indexOf(p);
  pesanan.splice(index, 1);
  quantityElement.textContent = 0;
  notesInput.value = "";
  console.log(pesanan);
});
