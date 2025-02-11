const express = require('express')
const bodyParser = require('body-parser')
const session = require('express-session')
const mongoose = require('mongoose')
const { GoogleGenerativeAI } = require('@google/generative-ai')
const authRoutes = require('./routes/authRoutes')
const mainRoutes = require('./routes/mainRoutes')
const landmarkRoutes = require('./routes/landmarkRoutes')
const apiRoutes = require('./routes/apiRoutes')
const quizQuestionRoutes = require('./routes/quizQuestionRoutes')

const authMiddleware = require('./middleware/authMiddleware')
require('dotenv').config()

const app = express()

mongoose
	.connect(process.env.MONGODB)
	.then(() => {
		console.info('Connected to MongoDB')
	})
	.catch(err => {
		console.error('Error: ', err)
	})

app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))

app.use(express.static('public'))

app.set('view engine', 'ejs')
app.set('views', __dirname + '/views')

app.use(
	session({
		secret: process.env.SECRET1,
		resave: false,
		saveUninitialized: true,
	})
)

// Route to handle Gemini API integration
const genAI = new GoogleGenerativeAI(process.env.GENAI)

app.get('/ai-generate', async (req, res) => {
	try {
		const prompt = req.query.prompt || 'Explain how AI works'
		console.log('Prompt received:', prompt)

		const model = genAI.getGenerativeModel({ model: 'gemini-1.5-flash' })
		const result = await model.generateContent(prompt)

		const generatedText =
			result.response.candidates[0]?.content?.parts[0]?.text ||
			'No content generated.'

		res.json({ success: true, response: generatedText })
	} catch (error) {
		console.error('Error in /ai-generate:', error)
		res
			.status(500)
			.json({ success: false, message: 'Failed to generate content.' })
	}
})

app.use(authMiddleware)
app.use('/auth', authRoutes)
app.use('/landmarks', landmarkRoutes)
app.use('/api', apiRoutes)
app.use('/quiz-questions', quizQuestionRoutes)

app.use('/', mainRoutes)

app.listen(3000, () => {
	console.log('Server is running on port 3000')
})
