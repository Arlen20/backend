module.exports = (req, res, next) => {
	res.locals.loggedIn = req.session.isLoggedIn || false
	res.locals.user = req.session.user || null
	next()
}
