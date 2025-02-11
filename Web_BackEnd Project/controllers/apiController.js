const axios = require('axios')

exports.stockMarketAPI = async (req, res) => {
	try {
		const { symbol } = req.query
		const response = await axios.get('https://www.alphavantage.co/query', {
			params: {
				function: 'TIME_SERIES_INTRADAY',
				symbol: symbol,
				interval: '5min',
				apikey: process.env.STOCK_MARKET_API,
			},
		})

		res.json(response.data)
	} catch (error) {
		console.error('Error fetching stock price:', error)
		res.status(500).json({ error: 'Internal server error' })
	}
}

exports.newsAPI = async (req, res) => {
	try {
		const apiUrl = 'https://newsapi.org/v2/top-headlines'

		const response = await axios.get(apiUrl, {
			params: {
				apiKey: process.env.NEWS_API,
				country: 'us',
				pageSize: 10,
				category: 'general',
			},
		})
		res.json(response.data)
	} catch (error) {
		console.error('Error fetching news articles:', error)
		res.status(500).json({ error: 'Internal server error' })
	}
}

exports.searchAPI = async (req, res) => {
	try {
		const { query } = req.query
		const apiUrl = 'https://api.gemini.com/v1/search' // Уточните базовый URL в документации GeminiAI

		const response = await axios.get(apiUrl, {
			params: {
				query: query,
				api_key: process.env.API_KEY, // Ваш ключ API
			},
		})

		res.json(response.data)
	} catch (error) {
		console.error('Error fetching data from GeminiAI:', error)
		res.status(500).json({ error: 'Internal server error' })
	}
}

console.log('Gemini API Key:', process.env.API_KEY)
