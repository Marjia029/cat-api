<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Browser</title>
    <link rel="stylesheet" href="https://unpkg.com/swiper/swiper-bundle.min.css" />
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="stylesheet" href="/static/css/voting.css">
    <link rel="stylesheet" href="/static/css/breed_search.css">
    <link rel="stylesheet" href="/static/css/favs.css"></link>
</head>
<body>
    <div class="content-container">
        <!-- Navigation -->
        <nav class="nav-tabs">
            <a href="#voting" class="nav-item" data-page="voting">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 20V4M4 12l8-8 8 8"/>
                </svg>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 4v16M4 12l8 8 8-8"/>
                </svg>
                Voting
            </a>
            <a href="#breeds" class="nav-item" data-page="breeds">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="M21 21l-4.35-4.35"/>
                </svg>
                Breeds
            </a>
            <a href="#favorites" class="nav-item" data-page="favorites">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78L12 21.23l8.84-8.84a5.5 5.5 0 0 0 0-7.78z"/>
                </svg>
                Favs
            </a>
        </nav>

        <!-- Content sections -->
        <div id="voting-content" class="page-content">
            <h1>Vote for a Cat!</h1>
            <div class="image-container">
                <img id="voting-image" src="">
            </div>
            <div class="buttons-container">
                <button class="heart-icon" onclick="handleFavorite()">❤</button>
                <div class="like-dislike-container">
                    <button class="like-button" onclick="handleVote('like')">👍</button>
                    <button class="dislike-button" onclick="handleVote('dislike')">👎</button>
                </div>
            </div>
        </div>

        <div id="breeds-content" class="page-content">
            <div class="search-container">
                <select id="breed-select">
                    <!-- Breeds will be populated by JavaScript -->
                </select>
            </div>
            <div class="swiper">
                <div class="swiper-wrapper">
                    <!-- Breed images will be populated by JavaScript -->
                </div>
                <div class="swiper-pagination"></div>
            </div>
            <div class="breed-info">
                <h2 class="breed-title"></h2>
                <p class="breed-description"></p>
                <a href="#" target="_blank" class="wiki-link">WIKIPEDIA</a>
            </div>
        </div>

        <div id="favorites-content" class="page-content">
            <h1>Your Favorite Cat Images</h1>
            <div class="favorites-container">
                <div id="favorites-controls">
                    <button id="grid-view-btn"><i class="fas fa-th"></i></button>
                    <button id="column-view-btn"><i class="fas fa-list"></i></button>
                </div>
                <ul id="favorites-list">
                    <!-- Favorites will be populated by JavaScript -->
                </ul>
            </div>
        </div>
    </div>

    <script src="https://unpkg.com/swiper/swiper-bundle.min.js"></script>
    <script src="/static/js/spa.js"></script>
    <script src="/static/js/fav_view.js"></script>
</body>
</html>