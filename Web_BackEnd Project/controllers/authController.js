const nodemailer = require('nodemailer')
const bcrypt = require('bcrypt')
const User = require('../models/userModel')
const crypto = require('crypto')
require('dotenv').config()

const transporter = nodemailer.createTransport({
	service: 'gmail',
	auth: {
		user: process.env.EMAIL_USER,
		pass: process.env.EMAIL_PASS,
	},
})

exports.registerUser = async (req, res) => {
	try {
		const { username, password, email, firstName, lastName, age, gender } =
			req.body

		const existingUsername = await User.findOne({ username })
		if (existingUsername) {
			return res.status(400).send('Username already exists')
		}

		const existingEmail = await User.findOne({ email })
		if (existingEmail) {
			return res.status(400).send('Email already registered')
		} else if (!email.endsWith('@gmail.com')) {
			return res.status(401).send('Invalid email')
		}

		const otp = crypto.randomInt(100000, 999999) // Generate a 6-digit OTP
		const otpExpires = new Date(Date.now() + 5 * 60 * 1000) // OTP expires in 5 minutes

		const hashedPassword = await bcrypt.hash(password, 10)

		const newUser = new User({
			username,
			password: hashedPassword,
			email,
			firstName,
			lastName,
			age,
			gender,
			role: 'regular',
			otp,
			otpExpires,
			isVerified: false,
		})

		await newUser.save()

		const mailOptions = {
			from: process.env.EMAIL_USER,
			to: email,
			subject: 'Verify Your OTP',
			text: `Your OTP is: ${otp}. It will expire in 5 minutes.`,
		}

		transporter.sendMail(mailOptions, (error, info) => {
			if (error) {
				console.log(error)
				return res.status(500).send('Error sending email')
			} else {
				console.log('Email sent: ' + info.response)
				res.status(200).send('Registration successful. Welcome email sent.')
			}
		})
	} catch (error) {
		console.error('Error:', error)
		res.status(500).send('Error registering user')
	}
}

exports.loginUser = async (req, res) => {
	const { username, password } = req.body
	try {
		const user = await User.findOne({ username })
		if (user && (await bcrypt.compare(password, user.password))) {
			req.session.user = user
			req.session.isLoggedIn = true
			res.redirect('/')
		} else {
			res.status(401).send('Invalid username or password')
		}
	} catch (error) {
		console.error('Login error:', error)
		res.status(500).send('Internal server error')
	}
}

exports.logoutUser = (req, res) => {
	req.session.destroy(err => {
		if (err) {
			return res.status(500).send('Error logging out')
		}
		res.redirect('/')
	})
}

exports.renderLoginForm = (req, res) => {
	res.render('login') // Ensure that you have a 'login.ejs' or 'login.html' file in your views directory
}

exports.renderRegisterForm = (req, res) => {
	res.render('register') // Ensure that you have a 'register.ejs' or 'register.html' file in your views directory
}

exports.verifyOtp = async (req, res) => {
	try {
		const { email, otp } = req.body

		const user = await User.findOne({ email })
		if (!user) {
			return res.status(400).send('User not found')
		}

		if (user.otpExpires < Date.now()) {
			return res.status(400).send('OTP has expired')
		}

		if (user.otp !== parseInt(otp, 10)) {
			return res.status(400).send('Invalid OTP')
		}

		user.isVerified = true
		await user.save()

		res.status(200).send('OTP verified successfully. You can now log in.')
	} catch (error) {
		console.error('Error verifying OTP:', error)
		res.status(500).send('Error verifying OTP')
	}
}
