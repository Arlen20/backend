const mongoose = require('mongoose')
const QuizQuestion = require('./models/quizQuestionModel')
require('dotenv').config()

// Connect to MongoDB
mongoose
	.connect(process.env.MONGODB, {
		useNewUrlParser: true,
		useUnifiedTopology: true,
	})
	.then(() => console.log('Connected to MongoDB'))
	.catch(err => console.error('Error connecting to MongoDB:', err))

// Define quiz questions
const quizQuestions = [
	{
		question: 'What countries does Kazakhstan border with? (choose 5)',
		options: [
			'Russia',
			'Mongolia',
			'China',
			'Kyrgyzstan',
			'Türkiye',
			'Tajikistan',
			'Turkmenistan',
			'Uzbekistan',
		],
		correctAnswer: [
			'Russia',
			'China',
			'Kyrgyzstan',
			'Turkmenistan',
			'Uzbekistan',
		],
	},
	{
		question: 'Who established Kazakh Khanate?',
		options: [
			'Janibek Khan and Kerei Khan',
			'Abu’l-Khayr',
			'Jochi Khan',
			'Kenesary Khan',
		],
		correctAnswer: ['Janibek Khan and Kerei Khan'],
	},
	{
		question:
			'How much did the Kazakh population grow during the reign of Kasym Khan?',
		options: ['1,000,000', '2,300,000', '700,000', '1,500,000'],
		correctAnswer: ['1,000,000'],
	},
	{
		question: 'Which Juz was the first to join the Russian Empire?',
		options: ['Junior Juz', 'Middle Juz', 'Senior Juz'],
		correctAnswer: ['Junior Juz'],
	},
	{
		question:
			'When did the Russian Empire officially abolish the Kazakh Khanate?',
		options: ['1847', '1854', '1874', '1881'],
		correctAnswer: ['1847'],
	},
	{
		question:
			"Which event corresponds to the statement 'This period of openness and reform contributed to the rise of nationalist movements in Soviet republics, including Kazakhstan. On December 16, 1991, Kazakhstan declared its independence following the dissolution of the Soviet Union.'?",
		options: [
			'Ethnic and Cultural Changes',
			'Nuclear Testing',
			'World War II',
			'Perestroika and Independence',
		],
		correctAnswer: ['Perestroika and Independence'],
	},
	// Add more questions as needed
]

// Insert quiz questions into the database
QuizQuestion.insertMany(quizQuestions)
	.then(() => {
		console.log('Quiz questions added successfully')
		mongoose.disconnect()
	})
	.catch(err => {
		console.error('Error adding quiz questions:', err)
		mongoose.disconnect()
	})
