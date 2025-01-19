<script lang="ts">
  export let currentPage: number;
  export let totalPages: number;
  export let totalItems: number;
  export let onPageChange: (page: number) => void;
  
  function handlePageChange(page: number) {
    if (page >= 1 && page <= totalPages) {
      onPageChange(page);
    }
  }
</script>

<div class="pagination">
  <button 
    disabled={currentPage === 1}
    on:click={() => handlePageChange(currentPage - 1)}
    aria-label="Previous page"
  >
    &lt;
  </button>
  
  <span>Page {currentPage} of {totalPages} ({totalItems} items)</span>
  
  <button 
    disabled={currentPage === totalPages}
    on:click={() => handlePageChange(currentPage + 1)}
    aria-label="Next page"
  >
    &gt;
  </button>
</div>

<style>
    .pagination {
      display: flex;
      align-items: center;
      gap: 1rem;
      width: fit-content;
      margin: 1rem 0;
      justify-content: flex-start;
    }
    
    button {
      padding: 0.5rem 1rem;
      cursor: pointer;
      min-width: 2.5rem;
    }
    
    button:disabled {
      cursor: not-allowed;
      opacity: 0.5;
    }
</style>
