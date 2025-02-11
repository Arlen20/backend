const mongoose = require('mongoose')

const quizQuestionSchema = new mongoose.Schema({
	question: { type: String, required: true },
	options: [{ type: String, required: true }],
	correctAnswer: [{ type: String, required: true }],
	createdDate: { type: Date, default: Date.now },
	updatedDate: { type: Date, default: Date.now },
})

module.exports = mongoose.model('QuizQuestion', quizQuestionSchema)
