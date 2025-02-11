const express = require('express')
const router = express.Router()
const mainController = require('../controllers/mainController')

router.get('/', mainController.indexPage)
router.get('/account', mainController.account)
router.get('/admin', mainController.adminPage)
router.get('/stockmarketapi', mainController.stockMarketAPI)
router.get('/newsapi', mainController.newsAPI)
router.get('/inRussianEmpire', mainController.inRussianEmpirePage)
router.get('/kazakhKhanate', mainController.kazakhKhanatePage)
router.get('/partOfUSSR', mainController.partOfUSSRPage)
router.get('/quiz', mainController.quizPage)
router.post('/quiz/submit', mainController.submitQuiz)
router.get('/quiz/questions', mainController.getQuizQuestions)

module.exports = router
