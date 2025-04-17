<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  // Modal state for image zoom
  let showImageModal = false;
  let modalImageUrl = '';
  let modalImageAlt = '';
  let zoomed = false;
  let dragging = false;
  let dragStart = { x: 0, y: 0 };
  let imgOffset = { x: 0, y: 0 };
  let imgCurrent = { x: 0, y: 0 };

  function openImageModal(url: string, alt: string) {
    modalImageUrl = url;
    modalImageAlt = alt;
    showImageModal = true;
    zoomed = false;
    imgOffset = { x: 0, y: 0 };
    imgCurrent = { x: 0, y: 0 };
  }
  function closeImageModal() {
    showImageModal = false;
    zoomed = false;
    dragging = false;
    imgOffset = { x: 0, y: 0 };
    imgCurrent = { x: 0, y: 0 };
  }
  function toggleZoom(e: MouseEvent) {
    zoomed = !zoomed;
    if (!zoomed) {
      imgOffset = { x: 0, y: 0 };
      imgCurrent = { x: 0, y: 0 };
    } else {
      // Center zoom on click point
      const img = e.currentTarget as HTMLImageElement;
      const rect = img.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      imgOffset = { x: -(x * 2 - rect.width) / 2, y: -(y * 2 - rect.height) / 2 };
      imgCurrent = { ...imgOffset };
    }
  }
  function startDrag(e: MouseEvent) {
    dragging = true;
    dragStart = { x: e.clientX, y: e.clientY };
  }
  function stopDrag() {
    dragging = false;
    imgOffset = { ...imgCurrent };
  }
  function onDrag(e: MouseEvent) {
    if (!dragging) return;
    const dx = e.clientX - dragStart.x;
    const dy = e.clientY - dragStart.y;
    imgCurrent = { x: imgOffset.x + dx, y: imgOffset.y + dy };
  }

  let merchItems = [];
  let movies = [];
  let showOrderSummary = false;
  let totalAmount = 0;
  let userName = '';
  let userEmail = '';

  const POSTER_PRICE = 10;
  const sizes = ['S', 'M', 'L', 'XL'];

  function handleCheckoutClick() {
    updateTotal();
    showOrderSummary = true;
    setTimeout(() => {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }, 0);
  }

  function getCurrentCart(): CartItem[] {
    const merchCart = merchItems
      .filter(item => {
        // If there are multiple inventory keys, require size selection
        if (Object.keys(item.inventory).length > 1) {
          return item.selectedSize && item.quantity && item.quantity > 0;
        }
        // If only one inventory key (e.g., stickers), just require quantity
        return item.quantity && item.quantity > 0;
      })
      .map(item => ({
        merchId: item.ID,
        name: item.Name,
        image_url: item.ImageURL,
        size: Object.keys(item.inventory).length > 1 ? item.selectedSize : Object.keys(item.inventory)[0] || undefined,
        quantity: item.quantity!,
        price: item.Price
      }));
    const movieCart = movies
      .filter(movie => movie.quantity && movie.quantity > 0)
      .map(movie => ({
        movieId: movie.ID,
        name: movie.Title,
        image_url: movie.PosterURL,
        quantity: movie.quantity!,
        price: POSTER_PRICE
      }));
    return [...merchCart, ...movieCart];
  }

  function updateTotal() {
    totalAmount = getCurrentCart().reduce((sum, item) => sum + (item.price * item.quantity), 0);
  }

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
        inventory: (item.sizes || []).reduce((acc: { [key: string]: number }, size: MerchSize) => {
          acc[size.Size] = size.Quantity;
          return acc;
        }, {}),
        selectedSize: '',
        quantity: 0
      }));
      const movieJson = await moviesResp.json();
      movies = (movieJson.data || []).map((movie: any) => ({ ...movie, quantity: 0 }));
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });

  function handleMerchChange() {
    updateTotal();
  }
  function handleMovieChange() {
    updateTotal();
  }

  async function submitOrder() {
    const items = getCurrentCart().map(item => ({
      merchandise_id: item.merchId || null,
      movie_id: item.movieId || null,
      size: item.size || null,
      quantity: item.quantity
    }));
    const orderData = {
      name: userName,
      email: userEmail,
      items
    };

    try {
      const response = await fetch('/api/order', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(orderData)
      });

      if (response.ok) {
        // Reset selections
        merchItems = merchItems.map(item => ({ ...item, selectedSize: '', quantity: 0 }));
        movies = movies.map(movie => ({ ...movie, quantity: 0 }));
        userName = '';
        userEmail = '';
        showOrderSummary = false;
        updateTotal();
        alert('Order placed successfully! See your email for payment and pickup details.');
      } else {
        const errorText = await response.text();
        alert('Failed to place order. Server says: ' + errorText);
        throw new Error('Failed to place order: ' + errorText);
      }
    } catch (error) {
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
      <img src={item.ImageURL} alt={item.Name} loading="lazy" class="clickable-image" on:click={() => openImageModal(item.ImageURL, item.Name)} />
      <div class="merch-title-row">
        <h3>{item.Name}</h3>
        <span class="price">${item.Price}</span>
      </div>
      <div class="merch-desc">{item.Description}</div>
      <div class="merch-controls">
        {#if Object.keys(item.inventory).length > 1}
          <select class="size-select" bind:value={item.selectedSize} on:change={handleMerchChange}>
  <option value="" class="select-size-placeholder">Select Size</option>
            {#each sizes.filter(size => item.inventory[size] !== undefined) as size}
  <option value={size} disabled={!item.inventory[size]}>
    {size} {item.inventory[size] > 0 ? '' : '(not available)'}
  </option>
{/each}
          </select>
          <input class="qty-input" type="number" min="0" max={item.selectedSize ? item.inventory[item.selectedSize] : 0} bind:value={item.quantity} on:input={handleMerchChange} placeholder="Qty" disabled={!item.selectedSize} />
        {:else}
          <input class="qty-input" type="number" min="0" max={item.inventory ? Object.values(item.inventory)[0] : 99} bind:value={item.quantity} on:input={handleMerchChange} placeholder="Qty" />
        {/if}
      </div>
    </div>
  {/each}
</div>
    </section>

    <section class="posters">
  <h2 class="poster-header">Movie Posters <span class="faq-tooltip-wrap"><span class="faq-icon" tabindex="0" aria-label="Movie posters FAQ">?</span><span class="faq-tooltip">Movie posters are 15 x 22 inches and printed to order on high quality paper.</span></span></h2>
  <div class="poster-price-line">$10 each</div>
  <div class="poster-grid">
    {#each movies as movie}
      <div class="poster-img-wrap">
        <img class="poster-img clickable-image" src={movie.PosterURL} alt={movie.Title} loading="lazy" on:click={() => openImageModal(movie.PosterURL, movie.Title)} />
        <div class="poster-qty-row">
          <input class="qty-input poster-qty" type="number" min="0" bind:value={movie.quantity} on:input={handleMovieChange} placeholder="Qty" />
        </div>
      </div>
    {/each}
  </div>
</section>

    <div class="checkout-section">
      <button 
        class="checkout-button" 
        on:click={handleCheckoutClick}
        disabled={getCurrentCart().length === 0}
      >
        Check Out (${totalAmount})
      </button>
    </div>
  {:else}
    <section class="order-summary">
      <h2>Order Summary</h2>
      {#if getCurrentCart().some(item => item.merchId)}
        <h3 class="summary-subheader">Merch</h3>
        <div class="summary-section">
          {#each getCurrentCart().filter(item => item.merchId) as item}
            <div class="summary-item">
              <div class="summary-img-wrap">
                <img src={item.image_url} alt={item.name} />
              </div>
              <div class="details">
                <div class="summary-title-wrap"><h4>{item.name}</h4></div>
                {#if item.size}
                  <p>Size: {item.size}</p>
                {/if}
                <p>Quantity: {item.quantity}</p>
                <p>Price: ${item.price * item.quantity}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
      {#if getCurrentCart().some(item => item.movieId)}
        <h3 class="summary-subheader">Posters</h3>
        <div class="summary-section">
          {#each getCurrentCart().filter(item => item.movieId) as item}
            <div class="summary-item">
              <div class="summary-img-wrap">
                <img src={item.image_url} alt={item.name} />
              </div>
              <div class="details">
                <div class="summary-title-wrap"><h4>{item.name}</h4></div>
                <p>Quantity: {item.quantity}</p>
                <p>Price: ${item.price * item.quantity}</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
      <div class="total">Total: ${totalAmount}</div>
      <div class="order-user-fields">
        <label>Name: <input type="text" bind:value={userName} required placeholder="Enter your name" /></label>
        <label>Email: <input type="email" bind:value={userEmail} required placeholder="Enter your email" /></label>
      </div>
      <div class="summary-controls">
        <button class="back-button" on:click={() => showOrderSummary = false}>Back to Shop</button>
        <button class="submit-button" on:click={submitOrder} disabled={!userName || !userEmail}>Submit Order</button>
      </div>
    </section>
  {/if}
{#if showImageModal}
  <div class="img-modal-backdrop" on:click={closeImageModal} role="button" tabindex="0" aria-label="Close image modal" on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeImageModal(); }}>
    <div class="img-modal-content" on:click|stopPropagation>
      <div class="img-modal-img-wrap">
        <span class="img-modal-x-gold"
          role="button"
          tabindex="0"
          aria-label="Close image"
          on:click={closeImageModal}
          on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeImageModal(); }}
        >&#10005;</span>
        <button
          type="button"
          class="img-modal-img-button"
          aria-label={zoomed ? 'Zoom out' : 'Zoom in'}
          style="padding:0; border:none; background:none; display:inline;"
          on:click={toggleZoom}
          on:mousedown={zoomed ? startDrag : undefined}
          on:mouseup={zoomed ? stopDrag : undefined}
          on:mouseleave={zoomed ? stopDrag : undefined}
          on:mousemove={zoomed && dragging ? onDrag : undefined}
        >
          <img
            src={modalImageUrl}
            alt={modalImageAlt}
            class="img-modal-img"
            style="cursor: {zoomed ? 'zoom-out' : 'zoom-in'}; transform: scale({zoomed ? 2 : 1}) translate({zoomed ? imgCurrent.x : 0}px, {zoomed ? imgCurrent.y : 0}px);"
            draggable="false"
          />
        </button>
      </div>
    </div>
  </div>
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

.poster-header {
  font-size: 2rem;
  margin-bottom: 1rem;
  display: inline-flex;
  align-items: center;
  gap: 1.1rem;
  text-align: center;
}
.posters {
  text-align: center;
}
.faq-tooltip-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
  height: 100%;
}
.faq-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25em;
  height: 1.25em;
  border-radius: 50%;
  background: #ffd700;
  color: #222;
  font-weight: bold;
  text-align: center;
  line-height: 1.25em;
  font-size: 0.5em;
  cursor: pointer;
  margin-left: 0.1em;
  margin-right: 0.1em;
  box-shadow: 0 1px 3px #0002;
  transition: background 0.2s;
}
.faq-icon:focus, .faq-icon:hover {
  background: #ffe066;
}
.faq-tooltip {
  visibility: hidden;
  opacity: 0;
  position: absolute;
  left: 50%;
  top: 120%;
  transform: translateX(-50%);
  background: #222;
  color: #ffd700;
  padding: 0.18em 0.5em;
  border-radius: 4px;
  font-size: 0.78rem;
  white-space: nowrap;
  box-shadow: 0 2px 8px #0005;
  z-index: 10;
  pointer-events: none;
  transition: opacity 0.18s;
}
@media screen and (max-width: 768px) {
  .faq-tooltip {
    min-width: 20vw;
    max-width: 120vw;
    white-space: normal;
    word-break: normal;
    left: 0;
    transform: none;
  }
}
.faq-tooltip-wrap:hover .faq-tooltip, .faq-tooltip-wrap:focus-within .faq-tooltip {
  visibility: visible;
  opacity: 1;
  pointer-events: auto;
}


  .merch-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 2rem;
  margin-bottom: 4rem;
}
.poster-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.13rem;
  margin-bottom: 4rem;
  justify-items: center;
}

  .merch-card {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 280px;
  box-sizing: border-box;
  gap: 0.75rem;
  transition: transform 0.2s;
}
.merch-card:hover {
  transform: translateY(-5px);
}
.poster-img-wrap {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.poster-qty-row {
  display: flex;
  justify-content: center;
  margin-top: 0.1rem;
}


img {
  width: 100%;
  height: auto;
  object-fit: contain;
  border-radius: 4px;
  background: #222;
}

.poster-img {
  width: 100%;
  max-width: 275px;
  height: auto;
  aspect-ratio: 2/3;
  object-fit: contain;
  border-radius: 4px;
  background: #222;
  margin-bottom: 0.25rem;
}
.clickable-image {
  cursor: zoom-in;
  transition: box-shadow 0.18s;
}
.clickable-image:hover {
  box-shadow: 0 0 0 2px #ffd70099;
}

/* Modal styles */
.img-modal-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.84);
  z-index: 1500;
  display: flex;
  align-items: center;
  justify-content: center;
}
.img-modal-content {
  position: relative;
  background: none;
  box-shadow: none;
  border-radius: 0;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.img-modal-img-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 92vw;
  max-height: 92vh;
  overflow: hidden;
  background: none;
}
.img-modal-img {
  max-width: 92vw;
  max-height: 92vh;
  user-select: none;
  border-radius: 0;
  background: none;
  transition: transform 0.16s;
  cursor: zoom-in;
}
.img-modal-img[style*="zoom-out"] {
  cursor: zoom-out;
}
.img-modal-img {
  /* Allow pan/zoom */
  will-change: transform;
}
.img-modal-x-gold {
  position: absolute;
  top: 0.5rem;
  right: 0.7rem;
  color: #ffd700;
  font-size: 1.7rem;
  font-family: inherit;
  font-weight: 700;
  background: none;
  border: none;
  cursor: pointer;
  z-index: 10;
  line-height: 1;
  padding: 0;
  user-select: none;
  transition: color 0.16s;
  outline: none;
}
.img-modal-x-gold:hover, .img-modal-x-gold:focus {
  color: #fffbe6;
  outline: 2px solid #ffd70044;
}
.img-modal-img-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 92vw;
  max-height: 92vh;
  overflow: hidden;
  background: none;
}



  h3 {
  font-size: 1.2rem;
  margin: 0;
  display: inline;
}
.merch-title-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 0.25rem;
}
.price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #ffd700;
}
.poster-price-line {
  color: #ffd700;
  font-size: 1.2rem;
  margin-top: 0.2rem;
  margin-bottom: 1.5rem;
  text-align: center;
}

  .merch-desc {
  font-size: 1rem;
  margin-bottom: 0.25rem;
  margin-top: 0.25rem;
  text-align: center;
  color: #eee;
}
.merch-controls {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}
.size-select {
  width: 110px;
  min-width: 110px;
  max-width: 110px;
}
.qty-input {
  width: 48px;
  min-width: 48px;
  max-width: 48px;
  text-align: center;
  font-size: 1rem;
}
.poster-qty {
  margin-bottom: 0.75rem;
}


  select, input {
  padding: 0.35rem 0.5rem;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  box-sizing: border-box;
}
select.size-select {
  font-size: 0.93rem;
  width: 110px;
  min-width: 110px;
  max-width: 110px;
}
option.select-size-placeholder {
  font-size: 0.9rem;
  color: #888;
}
input.qty-input {
  font-size: 1rem;
  width: 48px;
  min-width: 48px;
  max-width: 48px;
  text-align: center;
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
  }

  .checkout-button {
    font-size: 1.25rem;
    padding: 1rem 2rem;
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
  }

  .order-user-fields {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .order-user-fields label {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;
    color: #fff;
    width: auto;
  }
  .order-user-fields input {
    margin-top: 0;
    padding: 0.4rem 0.6rem;
    border-radius: 4px;
    border: none;
    font-size: 1rem;
    width: 180px;
    max-width: 60vw;
  }

  .order-summary {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .summary-subheader {
    font-size: 1.3rem;
    color: #ffd700;
    margin: 0.7rem 0 0.7rem 0;
    text-align: center;
    letter-spacing: 0.03em;
  }

  .summary-section {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .summary-item {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 2.5rem;
    margin-bottom: 1rem;
    background: rgba(255, 255, 255, 0.1);
    padding: 1.5rem 2.5rem;
    border-radius: 10px;
    width: 480px;
    max-width: 98vw;
    box-sizing: border-box;
    box-shadow: 0 2px 8px #0002;
    text-align: center;
  }

  .summary-img-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 140px;
    max-width: 180px;
    margin-right: 1.5rem;
  }
  .summary-img-wrap img {
    width: 140px;
    height: auto;
    object-fit: contain;
    border-radius: 6px;
    background: #222;
    box-shadow: 0 2px 8px #0003;
  }

  .details {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    align-items: center;
    min-width: 0;
    text-align: center;
  }
  .details p {
    text-align: center;
    margin-left: auto;
    margin-right: auto;
    width: 100%;
    margin-top: 0.18em;
    margin-bottom: 0.18em;
  }

  .summary-title-wrap {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 0.7rem;
    text-align: center;
  }
  .summary-title-wrap h4 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    text-align: center;
    width: 100%;
  }

  .total {
    font-size: 1.5rem;
    text-align: right;
    margin: 0.5rem 0rem 2rem 0rem;
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
  h2 {
    margin-top: 2.2rem;
  }
  .checkout-section {
    margin-top: 4rem;
  }
  /* Set border-box globally */
  *, *::before, *::after {
    box-sizing: border-box;
  }
  body, main.merch {
    overflow-x: hidden;
    margin: 0;
    padding: 0;
    width: 100%;
    max-width: 100vw;
  }
  main.merch {
    margin-top: 4rem;
    padding: 0;
    width: 100%;
    max-width: 100vw;
  }
  .merch-grid {
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    width: 100%;
    max-width: 100vw;
    margin: 0 auto;
    padding: 0;
  }
  .merch-card {
    width: 95vw;
    max-width: 360px;
    min-width: 0;
    box-sizing: border-box;
    margin: 0.5rem auto;
  }
  .poster-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 0.5rem;
    width: 92vw;
    max-width: 360px;
    margin: 0 auto;
    box-sizing: border-box;
  }
  .poster-img-wrap {
    max-width: 40vw;
    min-width: 0;
    box-sizing: border-box;
    margin: 0 auto;
  }
  .poster-img {
    max-width: 80vw;
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    margin: 0 auto;
  }
  .order-summary {
    width: 96vw;
    max-width: 380px;
    margin: 0.5rem auto;
    box-sizing: border-box;
    padding: 0.5rem 0.5rem;
  }
  .summary-item {
    flex-direction: column;
    align-items: stretch;
    width: 95vw;
    max-width: 360px;
    margin: 0.5rem auto;
    padding: 1rem 0.5rem;
    gap: 0.75rem;
    box-sizing: border-box;
  }
  .summary-img-wrap {
    margin: 0 auto 1rem auto;
    min-width: 0;
    max-width: 90vw;
    box-sizing: border-box;
  }
  .summary-img-wrap img {
    width: 90vw;
    max-width: 170px;
    height: auto;
    display: block;
    margin: 0 auto;
  }
}
</style>