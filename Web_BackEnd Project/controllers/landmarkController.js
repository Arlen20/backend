const Landmark = require('../models/landmarkModel');
const multer = require('multer');
const fs = require('fs');
const path = require('path');

const storage = multer.diskStorage({
    destination: function (req, file, cb) {
        cb(null, 'public/uploads/')
    },
    filename: function (req, file, cb) {
        cb(null, file.fieldname + '-' + Date.now() + path.extname(file.originalname))
    }
});

const upload = multer({ storage: storage });

exports.createLandmark = async (req, res) => {
    try {
        upload.array('pictures', 3)(req, res, async function (err) {
            if (err instanceof multer.MulterError) {
                return res.status(400).json({ message: 'File upload error' });
            } else if (err) {
                return res.status(400).json({ message: err.message });
            }

            const { title, description } = req.body;
            const pictures = req.files.map(file => file.path);

            const landmark = await Landmark.create({
                title,
                description,
                pictures
            });

            res.status(201).json(landmark);
        });
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
};

exports.getAllLandmarks = async (req, res) => {
    try {
        const landmarks = await Landmark.find();
        res.json(landmarks);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
};

exports.getLandmarkById = async (req, res) => {
    try {
        const landmark = await Landmark.findById(req.params.id);
        if (!landmark) {
            return res.status(404).json({ message: 'Landmark not found' });
        }
        res.json(landmark);
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
};

exports.updateLandmark = async (req, res) => {
    try {
        const updatedLandmark = await Landmark.findByIdAndUpdate(req.params.id, req.body, { new: true });
        if (!updatedLandmark) {
            return res.status(404).json({ message: 'Landmark not found' });
        }
        updatedLandmark.updatedDate = Date.now();
        await updatedLandmark.save();
        res.json(updatedLandmark);
    } catch (err) {
        res.status(400).json({ message: err.message });
    }
};

exports.deleteLandmark = async (req, res) => {
    try {
        const deletedLandmark = await Landmark.findByIdAndDelete(req.params.id);
        if (!deletedLandmark) {
            return res.status(404).json({ message: 'Landmark not found' });
        }

        deletedLandmark.pictures.forEach(picturePath => {
            console.log(picturePath);
            fs.unlinkSync(path.join(__dirname, '..', picturePath));
        });

        res.json({ message: 'Landmark deleted' });
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
};