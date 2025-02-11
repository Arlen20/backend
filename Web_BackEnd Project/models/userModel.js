const mongoose = require('mongoose')
const bcrypt = require('bcrypt')

const userSchema = new mongoose.Schema({
	username: { type: String, required: true, trim: true, unique: true },
	password: { type: String, required: true, trim: true },
	firstName: { type: String, required: true, trim: true },
	lastName: { type: String, required: true, trim: true },
	otp: { type: Number },
	otpExpires: { type: Date },
	isVerified: { type: Boolean, default: false },
	email: {
		type: String,
		required: true,
		trim: true,
		unique: true,
		validate: {
			validator: function (value) {
				return /\b[A-Za-z0-9._%+-]+@gmail\.com\b/.test(value)
			},
			message: props => `${props.value} is not a valid Gmail account!`,
		},
	},
	age: Number,
	gender: String,
	role: String,
})

userSchema.pre('save', async function (next) {
	const user = this
	if (!user.isModified('password')) return next()

	try {
		const hashedPassword = await bcrypt.hash(user.password, 10)
		user.password = hashedPassword
		next()
	} catch (error) {
		return next(error)
	}
})

module.exports = mongoose.model('User', userSchema)
