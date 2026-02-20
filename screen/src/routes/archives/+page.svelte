<script lang="ts">
  import { formatRuntime } from '$lib';
  import { onMount } from 'svelte';
  import { apiBaseUrl } from '$lib/api';

  type Movie = {
    Title: string;
    Date: string;
    Runtime: number;
    PosterURL: string;
    MenuURL: string;
  };

  type MonthGroup = {
    label: string;
    movies: Movie[];
  };

  let archive: Movie[] = [];
  let error = '';
  let loading = true;
  let selectedYear: number | null = null;
  let menuModal: { open: boolean; url: string; title: string } = { open: false, url: '', title: '' };

  function getYears(movies: Movie[]): number[] {
    const years = new Set(movies.map(m => new Date(m.Date).getFullYear()));
    return [...years].sort((a, b) => b - a);
  }

  function groupByMonth(movies: Movie[]): MonthGroup[] {
    const groups = new Map<string, Movie[]>();
    for (const movie of movies) {
      const d = new Date(movie.Date);
      const key = `${d.getFullYear()}-${String(d.getMonth()).padStart(2, '0')}`;
      const label = d.toLocaleString('en-US', { month: 'long', year: 'numeric' });
      if (!groups.has(key)) {
        groups.set(key, []);
      }
      groups.get(key)!.push(movie);
    }
    return [...groups.entries()]
      .sort(([a], [b]) => b.localeCompare(a))
      .map(([_, movies]) => ({
        label: new Date(movies[0].Date).toLocaleString('en-US', { month: 'long', year: 'numeric' }),
        movies,
      }));
  }

  function openMenu(url: string, title: string) {
    menuModal = { open: true, url, title };
  }

  function closeMenu() {
    menuModal = { open: false, url: '', title: '' };
  }

  $: years = getYears(archive);
  $: filtered = selectedYear
    ? archive.filter(m => new Date(m.Date).getFullYear() === selectedYear)
    : archive;
  $: grouped = groupByMonth(filtered);

  onMount(async () => {
    try {
      const response = await fetch(`${apiBaseUrl}/movie/archive`);
      const data = await response.json();

      if (data.success) {
        archive = data.data;
      } else {
        error = 'Failed to load the movie archive.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie archive.';
    } finally {
      loading = false;
    }
  });
</script>

<main class="archive">
  <section class="theater-info">
    <h1>The Golden Arm Archive</h1>
    <p class="subtitle">Discover our past screenings and crafted drinks.</p>
  </section>

  {#if error}
    <p class="error">{error}</p>
  {/if}

  {#if loading}
    <div class="loading">Loading...</div>
  {:else if archive.length > 0}
    <div class="year-filter">
      <button
        class="year-pill"
        class:active={selectedYear === null}
        on:click={() => selectedYear = null}
      >All</button>
      {#each years as year}
        <button
          class="year-pill"
          class:active={selectedYear === year}
          on:click={() => selectedYear = year}
        >{year}</button>
      {/each}
    </div>

    {#each grouped as group}
      <section class="month-section">
        <h2 class="month-header">{group.label}</h2>
        <div class="movie-grid">
          {#each group.movies as movie}
            <div class="movie-card">
              <img src={movie.PosterURL} alt="{movie.Title} poster" loading="lazy" />
              <div class="card-info">
                <span class="runtime">{formatRuntime(movie.Runtime)}</span>
                {#if movie.MenuURL}
                  <button class="view-menu" on:click={() => openMenu(movie.MenuURL, movie.Title)}>
                    View Menu
                  </button>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </section>
    {/each}
  {/if}
</main>

{#if menuModal.open}
  <div class="modal-overlay" on:click={closeMenu} on:keydown={(e) => e.key === 'Escape' && closeMenu()} role="button" tabindex="-1">
    <div class="menu-modal" on:click|stopPropagation role="dialog" aria-label="{menuModal.title} menu">
      <button class="modal-close" on:click={closeMenu}>&times;</button>
      <img src={menuModal.url} alt="{menuModal.title} menu" />
    </div>
  </div>
{/if}

<style>
  main.archive {
    padding: 1rem 2rem 2rem;
    color: #f0f0f0;
    min-height: 100vh;
    max-width: 1100px;
    margin-left: auto;
    margin-right: auto;
  }

  .theater-info {
    text-align: center;
    margin-bottom: 2rem;
  }

  h1 {
    text-align: center;
    font-size: 2rem;
    margin-bottom: 0.5rem;
  }

  .subtitle {
    color: #bbb;
    font-size: 1.1rem;
  }

  .loading {
    text-align: center;
    color: #888;
    padding: 4rem 0;
    font-size: 1.1rem;
  }

  .year-filter {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin-bottom: 2.5rem;
    flex-wrap: wrap;
  }

  .year-pill {
    background: transparent;
    border: 1px solid #444;
    color: #ccc;
    padding: 0.4rem 1.2rem;
    border-radius: 20px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s ease;
  }

  .year-pill:hover {
    border-color: var(--gold);
    color: var(--gold);
  }

  .year-pill.active {
    background: var(--gold);
    color: #000;
    border-color: var(--gold);
    font-weight: 600;
  }

  .month-section {
    margin-bottom: 2.5rem;
  }

  .month-header {
    font-size: 1.1rem;
    font-weight: 400;
    color: #888;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #333;
  }

  .movie-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.5rem;
  }

  .movie-card {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .movie-card img {
    width: 100%;
    height: auto;
    border-radius: 3px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
    transition: transform 0.2s ease;
  }

  .movie-card img:hover {
    transform: scale(1.03);
  }

  .card-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.3rem;
    margin-top: 0.6rem;
  }

  .runtime {
    font-size: 0.85rem;
    color: #999;
  }

  .view-menu {
    background: none;
    border: none;
    color: var(--gold);
    font-size: 0.85rem;
    cursor: pointer;
    padding: 0;
    transition: opacity 0.2s ease;
  }

  .view-menu:hover {
    opacity: 0.7;
  }

  /* Modal */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.85);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 2rem;
  }

  .menu-modal {
    position: relative;
    display: inline-block;
    max-width: 600px;
  }

  .menu-modal img {
    display: block;
    max-width: 100%;
    max-height: 85vh;
    height: auto;
    border-radius: 3px;
  }

  .modal-close {
    position: absolute;
    top: -2.25rem;
    right: -2.5rem;
    background: none;
    border: none;
    color: #fff;
    font-size: 1.6rem;
    cursor: pointer;
    line-height: 1;
  }

  .modal-close:hover {
    color: var(--gold);
  }

  .error {
    text-align: center;
    color: #ff6b6b;
  }

  @media (max-width: 768px) {
    h1 {
      margin-top: 5.5rem;
    }

    .movie-grid {
      grid-template-columns: 1fr;
      gap: 1.5rem;
      max-width: 300px;
      margin: 0 auto;
    }

    main.archive {
      padding: 1rem;
    }
  }
</style>
