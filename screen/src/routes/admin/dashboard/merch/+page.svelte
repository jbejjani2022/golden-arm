<script lang="ts">
    import { onMount } from 'svelte';
    import Pagination from '../Pagination.svelte';

    let merchandise: Array<any> = [];
    let error: string = '';

    // Fetch merchandise data on page load
    onMount(async () => {
        try {
            const response = await fetch('/api/merch/all');
            const data = await response.json();

            if (data.success) {
                // Convert sizes array to inventory object for each merchandise item
                merchandise = (data.data as MerchItemFromServer[]).map(item => {
                    if (!item.sizes) {
                        return {
                            ...item,
                            Inventory: {}
                        };
                    }
                    return {
                        ...item,
                        Inventory: item.sizes.reduce((acc: Record<string, number>, size: Size) => {
                            acc[size.Size] = size.Quantity;
                            return acc;
                        }, {})
                    };
                });
            } else {
                error = 'Failed to load merchandise data.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while fetching the merchandise data.';
        }
    });
    
    // Add Merchandise form data and handling
    let showAddForm = false;
    let newMerch = {
        Name: '',
        Description: '',
        Price: 0,
        ImageFile: null as File | null,
        Sizes: [
            { size: 'S', quantity: 0 },
            { size: 'M', quantity: 0 },
            { size: 'L', quantity: 0 }
        ]
    };

    const handleAddMerch = async () => {
        const formData = new FormData();
        formData.append('name', newMerch.Name);
        formData.append('description', newMerch.Description);
        formData.append('price', newMerch.Price.toString());
        if (newMerch.ImageFile) formData.append('image', newMerch.ImageFile);
        
        // Add sizes with quantities
        newMerch.Sizes.forEach(s => {
            formData.append('sizes', `${s.size},${s.quantity}`);
        });

        try {
            const response = await fetch('/api/merch', {
                method: 'POST',
                body: formData,
            });

            const result = await response.json();

            if (result.success) {
                // Update the local merchandise array instead of reloading
                const newResponse = await fetch('/api/merch/all');
                const newData = await newResponse.json();
                if (newData.success) {
                    showAddForm = false; // Close the form
                    merchandise = (newData.data as MerchItemFromServer[]).map(item => ({
                        ...item,
                        Inventory: item.sizes.reduce((acc: Record<string, number>, size: Size) => {
                            acc[size.Size] = size.Quantity;
                            return acc;
                        }, {})
                    }));
                }
            } else {
                error = result.error || 'Failed to add merchandise.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while adding the merchandise.';
        }
    };

    // Edit Merchandise form data and handling
    let showEditForm = false;
    interface Size {
        Size: string;
        Quantity: number;
        ID: string;
        MerchandiseID: string;
    }

    interface MerchItemFromServer {
        ID: string;
        Name: string;
        Description: string;
        Price: number;
        ImageURL: string;
        sizes: Size[];
    }

    interface MerchItem {
        ID: string;
        Name: string;
        Description: string;
        Price: number;
        Inventory: Record<string, number>;
        ImageURL: string;
        ImageFile?: File | null;
        Sizes: Array<{ size: string; quantity: number }>;
        // Original values for comparison
        originalName?: string;
        originalDescription?: string;
        originalPrice?: number;
        originalInventory?: Record<string, number>;
    }

    let editingMerch: MerchItem | null = null;
    
    const startEdit = (merch: MerchItem) => {
        editingMerch = {
            ...merch,
            ImageFile: null,
            // Store original values for comparison
            originalName: merch.Name,
            originalDescription: merch.Description,
            originalPrice: merch.Price,
            originalInventory: {...merch.Inventory},
            Sizes: Object.entries(merch.Inventory || {}).map(([size, quantity]) => ({
                size,
                quantity
            }))
        };
        showEditForm = true;
    };

    const handleEditMerch = async () => {
        if (!editingMerch) return;

        const formData = new FormData();
        // Only append fields that have been modified
        if (editingMerch.Name !== editingMerch.originalName) {
            formData.append('name', editingMerch.Name);
        }
        if (editingMerch.Description !== editingMerch.originalDescription) {
            formData.append('description', editingMerch.Description);
        }
        if (editingMerch.Price !== editingMerch.originalPrice) {
            formData.append('price', editingMerch.Price.toString());
        }
        if (editingMerch.ImageFile) {
            formData.append('image', editingMerch.ImageFile);
        }
        
        // Only append sizes that have changed
        editingMerch.Sizes.forEach(s => {
            const originalQuantity = editingMerch?.originalInventory?.[s.size] || 0;
            if (s.quantity !== originalQuantity) {
                formData.append('sizes', `${s.size},${s.quantity}`);
            }
        });

        try {
            const response = await fetch(`/api/merch/${editingMerch.ID}`, {
                method: 'PUT',
                body: formData,
            });

            const result = await response.json();
            if (result.success) {
                console.log('Merchandise updated successfully.');
                window.location.reload();
            } else {
                error = 'Failed to update merchandise.';
            }
        } catch (err) {
            console.error(err);
            error = 'Something went wrong while updating the merchandise.';
        }
    };

    // Delete merchandise handler
    const deleteMerch = async (merchId: string) => {
        const confirmation = confirm('Are you sure you want to delete this item?');

        if (confirmation) {
            try {
                const response = await fetch(`/api/merch/${merchId}`, {
                    method: 'DELETE',
                });

                if (response.ok) {
                    window.location.reload();
                } else {
                    console.error('Failed to delete merchandise');
                }
            } catch (err) {
                console.error('Error during merchandise deletion:', err);
            }
        }
    };

    // Pagination settings
    const ITEMS_PER_PAGE = 8;
    let currentPage = 1;

    $: paginatedMerch = merchandise.slice(
        (currentPage - 1) * ITEMS_PER_PAGE,
        currentPage * ITEMS_PER_PAGE
    );

    $: totalPages = Math.ceil(merchandise.length / ITEMS_PER_PAGE);

    const handlePageChange = (page: number) => {
        currentPage = page;
    };
</script>

<h2>Merchandise</h2>

<!-- Display error message if data fetching fails -->
{#if error}
    <p style="color: red;">{error}</p>
{/if}

<!-- Merchandise Table -->
{#if merchandise.length > 0}
<div class="table-container">
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Price</th>
                <th>Inventory</th>
                <th>Image</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {#each paginatedMerch as merch (merch.ID)}
                <tr>
                    <td>{merch.Name}</td>
                    <td>{merch.Description}</td>
                    <td>${merch.Price.toFixed(2)}</td>
                    <td>
                        {#if merch.Inventory && Object.keys(merch.Inventory).length > 0}
                            {#each Object.entries(merch.Inventory) as [size, quantity]}
                                {size}: {quantity}<br>
                            {/each}
                        {:else}
                            No inventory data
                        {/if}
                    </td>
                    <td>
                        <a href={merch.ImageURL} target="_blank">View</a>
                    </td>
                    <td>
                        <button class="link-button" on:click={() => startEdit(merch)}>Edit</button>
                        <button class="link-button delete" on:click={() => deleteMerch(merch.ID)}>Delete</button>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
    <div class="table-footer">
        <button class="add-button" on:click={() => showAddForm = !showAddForm}>Add Item</button>
        <div class="pagination-wrapper">
            <Pagination
                currentPage={currentPage}
                totalPages={totalPages}
                totalItems={merchandise.length}
                onPageChange={handlePageChange}
            />
        </div>
    </div>
</div>
{:else}
    <p>No merchandise found. Add some!</p>
    <button class="add-button" on:click={() => showAddForm = !showAddForm}>Add Item</button>
{/if}

<!-- Add Merchandise Form -->
{#if showAddForm}
<div class="modal">
    <div class="modal-content">
        <h2>Add Merchandise</h2>
        <form on:submit|preventDefault={handleAddMerch}>
            <div class="form-group">
                <label for="name">Name:</label>
                <input type="text" id="name" bind:value={newMerch.Name} required />
            </div>
            <div class="form-group">
                <label for="description">Description:</label>
                <textarea id="description" bind:value={newMerch.Description} required></textarea>
            </div>
            <div class="form-group">
                <label for="price">Price:</label>
                <input type="number" id="price" bind:value={newMerch.Price} required min="0" step="0.01"/>
            </div>
            <div class="form-group">
                <span id="inventory-label">Inventory:</span>
                {#each newMerch.Sizes as size}
                    <div class="size-input">
                        <label for="add-size-{size.size}">{size.size}:</label>
                        <input type="number" id="add-size-{size.size}" bind:value={size.quantity} min="0" required aria-labelledby="inventory-label" />
                    </div>
                {/each}
            </div>
            <div class="form-group">
                <label for="imageFile">Image:</label>
                <input type="file" id="imageFile" accept="image/*" on:change={(event) => newMerch.ImageFile = (event.target as HTMLInputElement).files?.[0] || null} required />
            </div>
            <button type="submit">Submit</button>
            <button type="button" class="cancel-button" on:click={() => showAddForm = false}>Cancel</button>
        </form>
    </div>
</div>
{/if}

<!-- Edit Merchandise Form -->
{#if showEditForm && editingMerch}
<div class="modal">
    <div class="modal-content">
        <h2>Edit Merchandise</h2>
        <p class="edit-note">Only modified fields will be updated.</p>
        <form on:submit|preventDefault={handleEditMerch}>
            <div class="form-group">
                <label for="edit-name">Name:</label>
                <input type="text" id="edit-name" bind:value={editingMerch.Name} />
            </div>
            <div class="form-group">
                <label for="edit-description">Description:</label>
                <textarea id="edit-description" bind:value={editingMerch.Description}></textarea>
            </div>
            <div class="form-group">
                <label for="edit-price">Price:</label>
                <input type="number" id="edit-price" bind:value={editingMerch.Price} min="0" step="0.01"/>
            </div>
            <div class="form-group">
                <span id="edit-inventory-label">Inventory:</span>
                {#each editingMerch.Sizes as size}
                    <div class="size-input">
                        <label for="edit-size-{size.size}">{size.size}:</label>
                        <input type="number" id="edit-size-{size.size}" bind:value={size.quantity} min="0" aria-labelledby="edit-inventory-label" />
                    </div>
                {/each}
            </div>
            <div class="form-group">
                <label for="edit-imageFile">New Image (optional):</label>
                <input type="file" id="edit-imageFile" accept="image/*" on:change={(event) => {
                    if (editingMerch) {
                        editingMerch.ImageFile = (event.target as HTMLInputElement).files?.[0] || null;
                    }
                }} />
            </div>
            <button type="submit">Update</button>
            <button type="button" class="cancel-button" on:click={() => showEditForm = false}>Cancel</button>
        </form>
    </div>
</div>
{/if}

<style>
    .table-container {
        width: 100%;
        margin-bottom: 2rem;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 20px;
    }

    th, td {
        padding: 12px;
        text-align: left;
        border: 1px solid #ddd;
    }

    th {
        background-color: #f4f4f4;
    }

    .table-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        margin-top: 1rem;
    }

    .pagination-wrapper {
        flex: 1;
        display: flex;
        justify-content: flex-end;
    }

    .modal {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-content {
        background: white;
        padding: 20px;
        border-radius: 5px;
        max-width: 500px;
        width: 90%;
    }

    .form-group {
        margin-bottom: 15px;
    }

    .form-group label {
        display: block;
        margin-bottom: 5px;
    }

    .form-group input[type="text"],
    .form-group input[type="number"],
    .form-group textarea {
        width: 100%;
        padding: 8px;
    }

    .link-button {
        background: none;
        border: none;
        padding: 0;
        color: #0066cc;
        cursor: pointer;
        font: inherit;
        margin-right: 1rem;
    }

    .link-button.delete {
        color: red;
    }

    .link-button:hover {
        text-decoration: underline;
    }

    .size-input {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 5px;
    }

    .size-input label {
        min-width: 30px;
    }

    .cancel-button {
        margin-left: 10px;
    }

    .edit-note {
        color: #666;
        font-style: italic;
        margin-bottom: 1rem;
    }

    .current-value {
        color: #666;
        font-size: 0.9em;
        font-weight: normal;
    }
</style>
