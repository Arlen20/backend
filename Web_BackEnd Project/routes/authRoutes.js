const express = require('express')
const router = express.Router()
const authController = require('../controllers/authController')

router.post('/register', authController.registerUser)
router.post('/login', authController.loginUser)
router.get('/logout', authController.logoutUser)
router.get('/login', authController.renderLoginForm)
router.get('/register', authController.renderRegisterForm)
router.post('/verify-otp', authController.verifyOtp)

module.exports = router
