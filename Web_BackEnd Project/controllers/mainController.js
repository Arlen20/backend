const Landmark = require('../models/landmarkModel')
const QuizResult = require('../models/QuizResult')
const Transaction = require('../models/transactionModel')
const User = require('../models/userModel')

exports.indexPage = async (req, res) => {
	const landmarks = await Landmark.find()
	landmarks.map(landmark => {
		landmark.pictures = landmark.pictures.map(picture => picture.substring(6))
		return landmark
	})
	const user = req.session.user
	const loggedIn = req.session.isLoggedIn
	let isAdmin = false
	if (loggedIn) {
		isAdmin = user.role === 'admin'
	}
	res.render('index', { user, isAdmin, loggedIn, landmarks })
}

exports.inRussianEmpirePage = async (req, res) => {
	res.render('inRussianEmpire')
}

exports.kazakhKhanatePage = async (req, res) => {
	res.render('kazakhKhanate')
}

exports.partOfUSSRPage = async (req, res) => {
	res.render('partOfUSSR')
}

exports.account = async (req, res) => {
	const user = req.session.user
	const loggedIn = req.session.isLoggedIn
	if (loggedIn) {
		let isAdmin = false
		isAdmin = user.role === 'admin'

		// Fetch the latest transaction for the user
		const transaction = await Transaction.findOne({
			'customer.id': user._id,
		}).sort({ updatedAt: -1 })

		res.render('account', { user, isAdmin, loggedIn, transaction })
	} else {
		res.redirect('/')
	}
}

exports.adminPage = (req, res) => {
	const user = req.session.user
	if (req.session.isLoggedIn) {
		if (user.role === 'admin') {
			res.render('admin', { user })
		} else {
			res.redirect('/')
		}
	} else {
		res.redirect('/')
	}
}

exports.stockMarketAPI = (req, res) => {
	const user = req.session.user
	const loggedIn = req.session.isLoggedIn
	if (loggedIn) {
		let isAdmin = false
		isAdmin = user.role === 'admin'
		res.render('stockmarketapi', { user, isAdmin, loggedIn })
	} else {
		res.redirect('/')
	}
}

exports.newsAPI = (req, res) => {
	const user = req.session.user
	const loggedIn = req.session.isLoggedIn
	if (loggedIn) {
		let isAdmin = false
		isAdmin = user.role === 'admin'
		res.render('newsapi', { user, isAdmin, loggedIn })
	} else {
		res.redirect('/')
	}
}

exports.submitQuiz = async (req, res) => {
	const { userId, answers } = req.body

	// Log the incoming request payload
	console.log('Received quiz submission:', req.body)

	// Basic validation
	if (!userId || !Array.isArray(answers) || answers.length === 0) {
		console.error('Missing required fields or invalid data format.')
		return res.status(400).json({
			success: false,
			message: 'Missing required fields or invalid data format.',
		})
	}

	try {
		let totalScore = 0

		// Validate answers
		for (const answer of answers) {
			const question = await QuizQuestion.findById(answer.questionId)
			const correctAnswers = question.correctAnswer.sort().join(', ')
			if (correctAnswers === answer.answer.split(', ').sort().join(', ')) {
				totalScore += 1
			}
		}

		totalScore = Math.round((totalScore / answers.length) * 100)

		const quizResult = new QuizResult({
			user: userId,
			totalScore,
			answers,
		})

		console.log('Saving quiz result:', quizResult)

		await quizResult.save()
		console.log('Quiz result saved successfully.')
		res.status(200).json({
			success: true,
			message: 'Quiz results saved successfully.',
			totalScore,
		})
	} catch (error) {
		console.error('Error saving quiz results:', error)
		res
			.status(500)
			.json({ success: false, message: 'Failed to save quiz results.' })
	}
}

exports.quizPage = async (req, res) => {
	try {
		const quizQuestions = await QuizQuestion.find()
		const user = req.session.user
		const loggedIn = req.session.isLoggedIn
		res.render('quiz', { loggedIn, user, quizQuestions })
	} catch (err) {
		console.error('Error rendering quiz page:', err)
		res.status(500).send('An error occurred while rendering quiz page.')
	}
}

exports.getQuizQuestions = async (req, res) => {
	try {
		const quizQuestions = await QuizQuestion.find()
		res.status(200).json({ quizQuestions })
	} catch (err) {
		console.error('Error fetching quiz questions:', err)
		res.status(500).send('An error occurred while fetching quiz questions.')
	}
}

exports.cartPage = (req, res) => {
	const user = req.session.user
	const loggedIn = req.session.isLoggedIn
	res.render('cart', { user, loggedIn })
}
