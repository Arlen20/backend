const express = require('express');
const router = express.Router();
const {createLandmark, getAllLandmarks, getLandmarkById, updateLandmark, deleteLandmark} = require('../controllers/landmarkController');

router.post('/', createLandmark);
router.get('/', getAllLandmarks);
router.get('/:id', getLandmarkById);
router.put('/:id', updateLandmark);
router.delete('/:id', deleteLandmark);

module.exports = router;