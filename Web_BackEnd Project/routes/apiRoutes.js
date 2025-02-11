const express = require('express')
const router = express.Router()
const { stockMarketAPI, newsAPI } = require('../controllers/apiController')

router.get('/stockmarketapi', stockMarketAPI)
router.get('/newsapi', newsAPI)

module.exports = router
