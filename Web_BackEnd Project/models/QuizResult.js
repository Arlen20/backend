const mongoose = require('mongoose')

const quizResultSchema = new mongoose.Schema({
	user: {
		type: mongoose.Schema.Types.ObjectId,
		ref: 'User',
		required: true,
	},
	totalScore: {
		type: Number,
		required: true,
	},
	answers: [
		{
			question: String,
			answer: String,
			correct: Boolean,
		},
	],
	createdAt: {
		type: Date,
		default: Date.now,
	},
})

const QuizResult = mongoose.model('QuizResult', quizResultSchema)
module.exports = QuizResult
