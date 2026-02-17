const mongoose = require('mongoose')

const cartItemSchema = new mongoose.Schema({
	id: String,
	name: String,
	price: Number,
	quantity: Number,
})

const customerSchema = new mongoose.Schema({
	id: String,
	name: String,
	email: String,
})

const transactionSchema = new mongoose.Schema({
	cartItems: [cartItemSchema],
	customer: customerSchema,
	status: String,
	createdAt: { type: Date, default: Date.now },
	updatedAt: { type: Date, default: Date.now },
	receipt_url: String,
})

const Transaction = mongoose.model('Transaction', transactionSchema)

module.exports = Transaction
