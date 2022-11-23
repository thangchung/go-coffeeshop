async function loadDatabase() {
    const db = await idb.openDB("tailwind_store", 1, {
      upgrade(db, oldVersion, newVersion, transaction) {
        db.createObjectStore("products", {
          keyPath: "id",
          autoIncrement: true,
        });
        db.createObjectStore("sales", {
          keyPath: "id",
          autoIncrement: true,
        });
      },
    });
  
    return {
      db,
      getProducts: async () => await db.getAll("products"),
      addProduct: async (product) => await db.add("products", product),
      editProduct: async (product) =>
        await db.put("products", product.id, product),
      deleteProduct: async (product) => await db.delete("products", product.id),
    };
  }
  
  function initApp() {
    const app = {
      db: null,
      time: null,
      firstTime: localStorage.getItem("first_time") === null,
      activeMenu: 'pos',
      loadingSampleData: false,
      moneys: [2000, 5000, 10000, 20000, 50000, 100000],
      products: [],
      keyword: "",
      cart: [],
      cash: 0,
      change: 0,
      isShowModalReceipt: false,
      receiptNo: null,
      receiptDate: null,
      async initDatabase() {
        this.db = await loadDatabase();
        this.loadProducts();
      },
      async loadProducts() {
        this.products = await this.db.getProducts();
        console.log("products loaded", this.products);
      },
      async startWithSampleData() {
        const response = await fetch("static/data/sample.json");
        const data = await response.json();
        this.products = data.products;
        for (let product of data.products) {
          await this.db.addProduct(product);
        }
  
        this.setFirstTime(false);
      },
      startBlank() {
        this.setFirstTime(false);
      },
      setFirstTime(firstTime) {
        this.firstTime = firstTime;
        if (firstTime) {
          localStorage.removeItem("first_time");
        } else {
          localStorage.setItem("first_time", new Date().getTime());
        }
      },
      filteredProducts() {
        const rg = this.keyword ? new RegExp(this.keyword, "gi") : null;
        return this.products.filter((p) => !rg || p.name.match(rg));
      },
      addToCart(product) {
        const index = this.findCartIndex(product);
        if (index === -1) {
          this.cart.push({
            productId: product.id,
            image: product.image,
            name: product.name,
            price: product.price,
            option: product.option,
            qty: 1,
          });
        } else {
          this.cart[index].qty += 1;
        }
        this.beep();
        this.updateChange();
      },
      findCartIndex(product) {
        return this.cart.findIndex((p) => p.productId === product.id);
      },
      addQty(item, qty) {
        const index = this.cart.findIndex((i) => i.productId === item.productId);
        if (index === -1) {
          return;
        }
        const afterAdd = item.qty + qty;
        if (afterAdd === 0) {
          this.cart.splice(index, 1);
          this.clearSound();
        } else {
          this.cart[index].qty = afterAdd;
          this.beep();
        }
        this.updateChange();
      },
      addCash(amount) {      
        this.cash = (this.cash || 0) + amount;
        this.updateChange();
        this.beep();
      },
      getItemsCount() {
        return this.cart.reduce((count, item) => count + item.qty, 0);
      },
      updateChange() {
        this.change = this.cash - this.getTotalPrice();
      },
      updateCash(value) {
        this.cash = parseFloat(value.replace(/[^0-9]+/g, ""));
        this.updateChange();
      },
      getTotalPrice() {
        return this.cart.reduce(
          (total, item) => total + item.qty * item.price,
          0
        );
      },
      submitable() {
        return this.change >= 0 && this.cart.length > 0;
      },
      submit() {
        const time = new Date();
        this.isShowModalReceipt = true;
        this.receiptNo = `TWPOS-KS-${Math.round(time.getTime() / 1000)}`;
        this.receiptDate = this.dateFormat(time);
      },
      closeModalReceipt() {
        this.isShowModalReceipt = false;
      },
      dateFormat(date) {
        const formatter = new Intl.DateTimeFormat('id', { dateStyle: 'short', timeStyle: 'short'});
        return formatter.format(date);
      },
      numberFormat(number) {
        return (number || "")
          .toString()
          .replace(/^0|\./g, "")
          .replace(/(\d)(?=(\d{3})+(?!\d))/g, "$1.");
      },
      priceFormat(number) {
        return number ? `Rp. ${this.numberFormat(number)}` : `Rp. 0`;
      },
      clear() {
        this.cash = 0;
        this.cart = [];
        this.receiptNo = null;
        this.receiptDate = null;
        this.updateChange();
        this.clearSound();
      },
      beep() {
        this.playSound("static/sound/beep-29.mp3");
      },
      clearSound() {
        this.playSound("static/sound/button-21.mp3");
      },
      playSound(src) {
        const sound = new Audio();
        sound.src = src;
        sound.play();
        sound.onended = () => delete(sound);
      },
      printAndProceed() {
        const receiptContent = document.getElementById('receipt-content');
        const titleBefore = document.title;
        const printArea = document.getElementById('print-area');
  
        printArea.innerHTML = receiptContent.innerHTML;
        document.title = this.receiptNo;
  
        window.print();
        this.isShowModalReceipt = false;
  
        printArea.innerHTML = '';
        document.title = titleBefore;
  
        // TODO save sale data to database
  
        this.clear();
      }
    };
  
    return app;
  }