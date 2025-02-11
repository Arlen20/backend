const express = require('express')
const router = express.Router()
const {
	createQuizQuestion,
	getAllQuizQuestions,
	getQuizQuestionById,
	updateQuizQuestion,
	deleteQuizQuestion,
} = require('../controllers/quizQuestionController')

router.post('/', createQuizQuestion)
router.get('/', getAllQuizQuestions)
router.get('/:id', getQuizQuestionById)
router.put('/:id', updateQuizQuestion)
router.delete('/:id', deleteQuizQuestion)

module.exports = router
