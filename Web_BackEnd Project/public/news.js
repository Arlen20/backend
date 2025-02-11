async function fetchNewsData() {
    try {
        const response = await fetch('/api/newsapi');
        const data = await response.json();
        return data.articles;
    } catch (error) {
        console.error('Error fetching news data:', error);
        return [];
    }
}

function displayNewsArticles(articles) {
    const newsContainer = document.getElementById('newsArticles');

    newsContainer.innerHTML = '';

    articles.forEach(article => {
        const articleElement = document.createElement('a');
        articleElement.classList.add('list-group-item', 'list-group-item-action');
        articleElement.href = article.url;
        articleElement.target = '_blank';
        articleElement.innerHTML = `
            <div class="d-flex w-100 justify-content-between">
                <h5 class="mb-1">${article.title}</h5>
                <small>${new Date(article.publishedAt).toDateString()}</small>
            </div>
            <p class="mb-1">${article.description}</p>
            <small>Source: ${article.source.name}</small>
        `;
        newsContainer.appendChild(articleElement);
    });
}

async function fetchAndDisplayNews() {
    const articles = await fetchNewsData();
    displayNewsArticles(articles);
}

document.addEventListener('DOMContentLoaded', fetchAndDisplayNews);