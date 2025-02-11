var stockChart = null;

function fetchStockData(companySymbol) {
    document.getElementById('title').innerHTML = 'Stock Price of ';
    if (companySymbol === 'AAPL'){
        document.getElementById('title').innerHTML += 'Apple';
    } else if (companySymbol === 'MSFT'){
        document.getElementById('title').innerHTML += 'Microsoft'
    } else if (companySymbol === 'GOOGL'){
        document.getElementById('title').innerHTML += 'Google'
    }
    fetch(`/api/stockmarketapi?symbol=${companySymbol}`)
        .then(response => response.json())
        .then(data => {
            const canvas = document.getElementById('stockChart');
            const ctx = canvas.getContext('2d');
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            const dates = Object.keys(data['Time Series (5min)']);
            const prices = dates.map(date => parseFloat(data['Time Series (5min)'][date]['4. close']));

            if (stockChart) {
                stockChart.destroy();
            }

            stockChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: dates,
                    datasets: [{
                        label: 'Stock Price (USD)',
                        data: prices,
                        borderColor: 'blue',
                        backgroundColor: 'rgba(0, 0, 255, 0.1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'day'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Price (USD)'
                            }
                        }
                    }
                }
            });
        })
        .catch(error => {
            console.error('Error fetching stock price data:', error);
        });
}

fetchStockData("AAPL")