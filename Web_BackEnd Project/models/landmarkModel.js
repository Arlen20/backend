const mongoose = require('mongoose');

const landmarkSchema = new mongoose.Schema({
    title: { type: String, required: true },
    description: { type: String, required: true },
    pictures: [{ type: String, required: true, max: 3 }],
    publishedDate: { type: Date, default: Date.now },
    updatedDate: { type: Date, default: Date.now }
});

module.exports = mongoose.model('Landmark', landmarkSchema);