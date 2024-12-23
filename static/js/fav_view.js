document.addEventListener("DOMContentLoaded", () => {
    const gridViewBtn = document.getElementById("grid-view-btn");
    const columnViewBtn = document.getElementById("column-view-btn");
    const favoritesList = document.getElementById("favorites-list");
  
    if (!gridViewBtn || !columnViewBtn || !favoritesList) {
      console.error("One or more elements not found. Ensure IDs are correct.");
      return;
    }
  
    // Ensure the favorites list starts with grid-view
    favoritesList.classList.add("grid-view");
  
    // Event listener for Grid View button
    gridViewBtn.addEventListener("click", () => {
      favoritesList.classList.remove("column-view");
      favoritesList.classList.add("grid-view");
      console.log("Switched to grid view");
    });
  
    // Event listener for Column View button
    columnViewBtn.addEventListener("click", () => {
      favoritesList.classList.remove("grid-view");
      favoritesList.classList.add("column-view");
      console.log("Switched to column view");
    });
  });
  