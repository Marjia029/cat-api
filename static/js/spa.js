// State management
let currentPage = 'voting';
let swiper = null;
let currentImageUrl = '';
let favorites = [];
let currentBreedId = null;

// API endpoints
const API_KEY = 'live_GWXcPdnWze27MNMJSjinKshtfsnVsi4EdrXfKUNhOmXsLakl5N7MwJCShLvC5Rxo';
// const API_ENDPOINTS = {
//     randomCat: 'https://api.thecatapi.com/v1/images/search',
//     breeds: 'https://api.thecatapi.com/v1/breeds',
//     breedImages: (breedId) => `https://api.thecatapi.com/v1/images/search?breed_ids=${breedId}&limit=8`
// };

// Initialize application
document.addEventListener('DOMContentLoaded', function() {
    setupNavigation();
    setupBreedSelect();

    // Load favorites from localStorage
    const savedFavorites = localStorage.getItem('favorites');
    if (savedFavorites) {
        favorites = JSON.parse(savedFavorites);
    }

    loadInitialPage();
});

// Navigation setup
function setupNavigation() {
    document.querySelectorAll('.nav-item').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const page = e.currentTarget.getAttribute('data-page');
            navigateToPage(page);
        });
    });
}

// Update page navigation function
async function navigateToPage(page) {
    // Hide all content sections
    document.querySelectorAll('.page-content').forEach(content => {
        content.style.display = 'none';
    });

    // If navigating to voting page, clear current image while loading
    if (page === 'voting') {
        const votingImage = document.getElementById('voting-image');
        votingImage.src = '';
        currentImageUrl = '';
    }

    // Show selected content
    document.getElementById(`${page}-content`).style.display = 'block';
    
    // Update active tab
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-page') === page) {
            item.classList.add('active');
        }
    });

    // Load page-specific content
    switch(page) {
        case 'voting':
            await loadRandomCat();
            break;
        case 'breeds':
            if (!document.querySelector('#breed-select').options.length) {
                await loadBreeds();
            } else if (currentBreedId) {
                // Reload current breed details when returning to breeds page
                await loadBreedDetails(currentBreedId);
            }
            break;
        case 'favorites':
            displayFavorites();
            break;
    }

    currentPage = page;
}

// Updated loadRandomCat function to return a promise
async function loadRandomCat() {
    try {
        // Fetch random cat image from VotingController (GET method)
        const response = await fetch('/voting', { method: 'GET' });
        const data = await response.json();
        currentImageUrl = data.image_url;
        document.getElementById('voting-image').src = data.image_url;
    } catch (error) {
        console.error('Error loading random cat:', error);
    }
}


async function handleVote(action) {
    try {
        const response = await fetch('/voting', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `action=${action}&image_url=${currentImageUrl}`
        });
        const data = await response.json();
        currentImageUrl = data.image_url;
        document.getElementById('voting-image').src = data.image_url;
        favorites = data.favorites;
    } catch (error) {
        console.error('Error handling vote:', error);
    }
}

async function handleFavorite() {
    try {
        const response = await fetch('/voting', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `action=favorite&image_url=${currentImageUrl}`
        });
        const data = await response.json();
        
        currentImageUrl = data.image_url;
        document.getElementById('voting-image').src = data.image_url;
        
        // Update favorites array and save to localStorage
        favorites = data.favorites;
        localStorage.setItem('favorites', JSON.stringify(favorites));

        if (currentPage === 'favorites') {
            displayFavorites();
        }
    } catch (error) {
        console.error('Error handling favorite:', error);
    }
}


// Update setupBreedSelect function
async function setupBreedSelect() {
    const select = document.querySelector('#breed-select');
    select.addEventListener('change', async (e) => {
        const breedId = e.target.value;
        await loadBreedDetails(breedId);
    });
}

async function loadBreeds() {
    try {
        // Fetch breed list from the Go server's BreedSearchController (GET method)
        const response = await fetch('/breed-search', { method: 'GET' });
        if (!response.ok) {
            throw new Error(`Failed to fetch breeds: ${response.statusText}`);
        }

        const breeds = await response.json();

        const select = document.querySelector('#breed-select');
        select.innerHTML = breeds.map(breed => 
            `<option value="${breed.id}">${breed.name}</option>`
        ).join('');

        // Automatically load details of the first breed
        if (breeds.length > 0) {
            await loadBreedDetails(breeds[0].id);
        }
    } catch (error) {
        console.error('Error loading breeds:', error);
    }
}

// Update loadBreedDetails function
async function loadBreedDetails(breedId) {
    try {
        const response = await fetch('/breed-search', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `breed_id=${breedId}`
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch breed details: ${response.statusText}`);
        }

        const data = await response.json();
        const selectedBreed = data.breed;
        const breedImages = data.images;

        // Update breed details
        document.querySelector('.breed-title').innerHTML = 
            `${selectedBreed.name} ${selectedBreed.origin ? `(${selectedBreed.origin})` : ''} <span class="breed-id">${selectedBreed.id}</span>`;
        document.querySelector('.breed-description').textContent = selectedBreed.description;

        const wikiLink = document.querySelector('.wiki-link');
        if (selectedBreed.wikipedia_url) {
            wikiLink.href = selectedBreed.wikipedia_url;
            wikiLink.style.display = 'inline';
        } else {
            wikiLink.style.display = 'none';
        }

        // Update breed images in the Swiper carousel
        const swiperWrapper = document.querySelector('.swiper-wrapper');
        swiperWrapper.innerHTML = breedImages.map(image => 
            `<div class="swiper-slide"><img src="${image.url}" alt="Breed Image"></div>`
        ).join('');

        // Reinitialize Swiper with updated slides
        if (swiper) {
            swiper.destroy(true, true);
        }

        swiper = new Swiper('.swiper', {
            slidesPerView: 1,
            spaceBetween: 0,
            loop: true,
            autoplay: {
                delay: 3000,
                disableOnInteraction: false,
                pauseOnMouseEnter: true,
            },
            pagination: {
                el: '.swiper-pagination',
                clickable: true,
            },
        });

    } catch (error) {
        console.error('Error loading breed details:', error);
    }
}

// Favorites page functions
function displayFavorites() {
    const favoritesList = document.getElementById('favorites-list');
    const reversedFavorites = [...favorites].reverse();
    if (favorites.length > 0) {
        favoritesList.innerHTML = reversedFavorites.map(url => 
            `<li><img src="${url}" alt="Favorite Cat Image" width="200"></li>`
        ).join('');
    } else {
        favoritesList.innerHTML = '<p>You have no favorite cat images yet.</p>';
    }
}

// Load initial page based on URL hash or default to voting
function loadInitialPage() {
    const hash = window.location.hash.slice(1);
    navigateToPage(hash || 'voting');
}

// Handle browser back/forward buttons
window.addEventListener('popstate', () => {
    loadInitialPage();
});