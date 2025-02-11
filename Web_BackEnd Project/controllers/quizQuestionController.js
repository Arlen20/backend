const QuizQuestion = require('../models/quizQuestionModel')

exports.createQuizQuestion = async (req, res) => {
	try {
		const { question, options, correctAnswer } = req.body
		const quizQuestion = await QuizQuestion.create({
			question,
			options,
			correctAnswer,
		})
		res.status(201).json(quizQuestion)
	} catch (err) {
		res.status(500).json({ message: err.message })
	}
}

exports.getAllQuizQuestions = async (req, res) => {
	try {
		const quizQuestions = await QuizQuestion.find()
		res.json(quizQuestions)
	} catch (err) {
		res.status(500).json({ message: err.message })
	}
}

exports.getQuizQuestionById = async (req, res) => {
	try {
		const quizQuestion = await QuizQuestion.findById(req.params.id)
		if (!quizQuestion) {
			return res.status(404).json({ message: 'Quiz Question not found' })
		}
		res.json(quizQuestion)
	} catch (err) {
		res.status(500).json({ message: err.message })
	}
}

exports.updateQuizQuestion = async (req, res) => {
	try {
		const updatedQuizQuestion = await QuizQuestion.findByIdAndUpdate(
			req.params.id,
			req.body,
			{ new: true }
		)
		if (!updatedQuizQuestion) {
			return res.status(404).json({ message: 'Quiz Question not found' })
		}
		updatedQuizQuestion.updatedDate = Date.now()
		await updatedQuizQuestion.save()
		res.json(updatedQuizQuestion)
	} catch (err) {
		res.status(400).json({ message: err.message })
	}
}

exports.deleteQuizQuestion = async (req, res) => {
	try {
		const deletedQuizQuestion = await QuizQuestion.findByIdAndDelete(
			req.params.id
		)
		if (!deletedQuizQuestion) {
			return res.status(404).json({ message: 'Quiz Question not found' })
		}
		res.json({ message: 'Quiz Question deleted' })
	} catch (err) {
		res.status(500).json({ message: err.message })
	}
}
