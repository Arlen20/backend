const express = require('express')
const router = express.Router()
const {
	indexPage,
	account,
	adminPage,
	stockMarketAPI,
	newsAPI,
	inRussianEmpirePage,
	kazakhKhanatePage,
	partOfUSSRPage,
	quizPage,
	cartPage,
} = require('../controllers/mainController')

router.get('/', indexPage)
router.get('/account', account)
router.get('/admin', adminPage)
router.get('/stockmarketapi', stockMarketAPI)
router.get('/newsapi', newsAPI)
router.get('/inRussianEmpire', inRussianEmpirePage)
router.get('/kazakhKhanate', kazakhKhanatePage)
router.get('/partOfUSSR', partOfUSSRPage)
router.get('/quiz', quizPage)
router.get('/cart', cartPage)

module.exports = router
