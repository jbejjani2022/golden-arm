<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  interface MerchSize {
    ID: string;
    MerchandiseID: string;
    Size: string;
    Quantity: number;
  }

  interface MerchItem {
    ID: string;
    Name: string;
    Description: string;
    ImageURL: string;
    Price: number;
    Sizes: MerchSize[];
    inventory: { [key: string]: number };
    selectedSize?: string;
    quantity?: number;
  }

  interface Movie {
    ID: string;
    Title: string;
    PosterURL: string;
    quantity?: number;
  }

  interface CartItem {
    merchId?: string;
    movieId?: string;
    name: string;
    image_url: string;
    size?: string;
    quantity: number;
    price: number;
  }

  let merchItems: MerchItem[] = [];
  let movies: Movie[] = [];
  let cart: CartItem[] = [];
  let showOrderSummary = false;
  let totalAmount = 0;

  const POSTER_PRICE = 10;
  const sizes = ['S', 'M', 'L', 'XL'];

  onMount(async () => {
    try {
      const [merchResp, moviesResp] = await Promise.all([
        fetch('/api/merch/all'),
        fetch('/api/movie/all')
      ]);
      const merchJson = await merchResp.json();
      console.log('Merch API Response:', merchJson);
      const merchData = merchJson.data;
      if (!Array.isArray(merchData)) {
        console.error('Expected data to be an array but got:', typeof merchData);
        merchItems = [];
        return;
      }
      merchItems = merchData.map((item: any) => ({
        ...item,
        inventory: item.sizes.reduce((acc: { [key: string]: number }, size: MerchSize) => {
          acc[size.Size] = size.Quantity;
          return acc;
        }, {})
      }));
      const movieJson = await moviesResp.json();
      movies = movieJson.data || [];
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });

  function addMerchToCart(item: MerchItem, size: string, quantity: number) {
    if (quantity <= 0 || !item.inventory[size]) return;

    const existingItem = cart.find(i => i.merchId === item.ID && i.size === size);
    if (existingItem) {
      existingItem.quantity = quantity;
      cart = [...cart];
    } else {
      cart = [...cart, {
        merchId: item.ID,
        name: item.Name,
        image_url: item.ImageURL,
        size,
        quantity,
        price: item.Price
      }];
    }
    updateTotal();
  }

  function addPosterToCart(movie: Movie, quantity: number) {
    if (quantity <= 0) return;

    const existingItem = cart.find(i => i.movieId === movie.ID);
    if (existingItem) {
      existingItem.quantity = quantity;
      cart = [...cart];
    } else {
      cart = [...cart, {
        movieId: movie.ID,
        name: movie.Title,
        image_url: movie.PosterURL,
        quantity,
        price: POSTER_PRICE
      }];
    }
    updateTotal();
  }

  function updateTotal() {
    totalAmount = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  }

  async function submitOrder() {
    const orderData = {
      items: cart.map(item => ({
        merch_id: item.merchId || null,
        movie_id: item.movieId || null,
        size: item.size || null,
        quantity: item.quantity
      }))
    };

    try {
      const response = await fetch('/api/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(orderData)
      });

      if (response.ok) {
        cart = [];
        showOrderSummary = false;
        alert('Order placed successfully!');
      } else {
        throw new Error('Failed to place order');
      }
    } catch (error) {
      console.error('Error submitting order:', error);
      alert('Failed to place order. Please try again.');
    }
  }
</script>

<main class="merch">
  {#if !showOrderSummary}
    <section class="merch-items">
      <h2>The Golden Arm Shop</h2>
      <div class="merch-grid">
        {#each merchItems as item}
          <div class="merch-card">
            <img src={item.ImageURL} alt={item.Name} />
            <h3>{item.Name}</h3>
            <p class="description">{item.Description}</p>
            <p class="price">${item.Price}</p>
            
            <div class="controls">
              <select bind:value={item.selectedSize}>
                <option value="">Select Size</option>
                {#each sizes as size}
                  <option value={size} disabled={!item.inventory[size]}>
                    {size} {item.inventory[size] > 0 ? '' : '(not available)'}
                  </option>
                {/each}
              </select>
              
              <input 
                type="number" 
                bind:value={item.quantity} 
                min="0" 
                max={item.selectedSize ? item.inventory[item.selectedSize] : 0}
              />
              
              <button 
                on:click={() => item.selectedSize && item.quantity && addMerchToCart(item, item.selectedSize, item.quantity)}
                disabled={!item.selectedSize || !item.quantity}
              >
                Add to Cart
              </button>
            </div>
          </div>
        {/each}
      </div>
    </section>

    <section class="posters">
      <h2>Movie Posters (${POSTER_PRICE} each)</h2>
      <div class="poster-grid">
        {#each movies as movie}
          <div class="poster-card">
            <img src={movie.PosterURL} alt={movie.Title} />
            <h3>{movie.Title}</h3>
            <div class="controls">
              <input 
                type="number" 
                bind:value={movie.quantity} 
                min="0"
              />
              <button 
                on:click={() => movie.quantity && addPosterToCart(movie, movie.quantity)}
                disabled={!movie.quantity}
              >
                Add to Cart
              </button>
            </div>
          </div>
        {/each}
      </div>
    </section>

    <div class="checkout-section">
      <button 
        class="checkout-button" 
        on:click={() => showOrderSummary = true}
        disabled={cart.length === 0}
      >
        Check Out (${totalAmount})
      </button>
    </div>
  {:else}
    <section class="order-summary">
      <h2>Order Summary</h2>
      
      {#if cart.some(item => item.merchId)}
        <div class="summary-section">
          {#each cart.filter(item => item.merchId) as item}
            <div class="summary-item">
              <img src={item.image_url} alt={item.name} />
              <div class="details">
                <h4>{item.name}</h4>
                <p>Size: {item.size}</p>
                <p>Quantity: {item.quantity}</p>
                <p>Price: ${item.price * item.quantity}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      {#if cart.some(item => item.movieId)}
        <div class="summary-section">
          {#each cart.filter(item => item.movieId) as item}
            <div class="summary-item">
              <img src={item.image_url} alt={item.name} />
              <div class="details">
                <h4>{item.name}</h4>
                <p>Quantity: {item.quantity}</p>
                <p>Price: ${item.price * item.quantity}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      <div class="total">Total: ${totalAmount}</div>
      
      <div class="summary-controls">
        <button class="back-button" on:click={() => showOrderSummary = false}>Back to Shop</button>
        <button class="submit-button" on:click={submitOrder}>Submit Order</button>
      </div>
    </section>
  {/if}
</main>

<style>
  main.merch {
    padding: 2rem;
    color: #ffffff;
    max-width: 1200px;
    margin: 0 auto;
  }

  h2 {
    font-size: 2rem;
    margin-bottom: 2rem;
    text-align: center;
  }

  .merch-grid, .poster-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 2rem;
    margin-bottom: 4rem;
  }

  .merch-card, .poster-card {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    transition: transform 0.2s;
  }

  .merch-card:hover, .poster-card:hover {
    transform: translateY(-5px);
  }

  img {
    width: 100%;
    height: 300px;
    object-fit: cover;
    border-radius: 4px;
  }

  h3 {
    font-size: 1.5rem;
    margin: 0;
  }

  .description {
    flex-grow: 1;
  }

  .price {
    font-size: 1.25rem;
    font-weight: bold;
    color: #ffd700;
  }

  .controls {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  select, input {
    padding: 0.5rem;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.9);
    border: none;
  }

  button {
    padding: 0.5rem 1rem;
    border-radius: 4px;
    border: none;
    background: #ffd700;
    color: black;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .checkout-section {
    position: sticky;
    bottom: 2rem;
    display: flex;
    justify-content: center;
    padding: 1rem;
  }

  .checkout-button {
    font-size: 1.25rem;
    padding: 1rem 2rem;
  }

  .order-summary {
    max-width: 800px;
    margin: 0 auto;
  }

  .summary-section {
    margin-bottom: 2rem;
  }

  .summary-item {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    background: rgba(255, 255, 255, 0.1);
    padding: 1rem;
    border-radius: 4px;
  }

  .summary-item img {
    width: 100px;
    height: 100px;
  }

  .details {
    flex-grow: 1;
  }

  .total {
    font-size: 1.5rem;
    text-align: right;
    margin: 2rem 0;
    color: #ffd700;
  }

  .summary-controls {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }

  .back-button {
    background: #666;
  }

  @media screen and (max-width: 768px) {
    main.merch {
      margin-top: 4rem;
      padding: 1rem;
    }

    .merch-grid, .poster-grid {
      grid-template-columns: 1fr;
    }

    .controls {
      flex-direction: column;
    }

    .summary-item {
      flex-direction: column;
    }

    .summary-item img {
      width: 100%;
      height: 200px;
    }
  }
</style>