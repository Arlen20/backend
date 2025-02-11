document
	.getElementById('registerForm')
	.addEventListener('submit', async event => {
		event.preventDefault()

		const username = document.getElementById('username').value.trim()
		const email = document.getElementById('email').value.trim()
		const password = document.getElementById('password').value.trim()
		const firstName = document.getElementById('firstName').value.trim()
		const lastName = document.getElementById('lastName').value.trim()
		const age = document.getElementById('age').value
		const gender = document.getElementById('gender').value

		if (
			username === '' ||
			email === '' ||
			password === '' ||
			firstName === '' ||
			lastName === ''
		) {
			displayMessage('Please fill in all fields.')
			return
		}

		const formData = {
			username: username,
			email: email,
			password: password,
			firstName: firstName,
			lastName: lastName,
			age: age,
			gender: gender,
		}

		try {
			const response = await fetch('/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(formData),
			})

			if (response.ok) {
				const otpInput = prompt('Enter the OTP sent to your email')
				if (otpInput) {
					const otpResponse = await fetch('/auth/verify-otp', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ email, otp: otpInput }),
					})
					if (otpResponse.ok) {
						window.location.href = '/'
					} else {
						const errorMessage = await otpResponse.text()
						console.error('OTP verification failed:', errorMessage)
						displayError(errorMessage)
					}
				}
			} else {
				const errorMessage = await response.text()
				console.error('Registration failed:', errorMessage)
				displayError(errorMessage)
			}
		} catch (error) {
			console.error('Error:', error)
		}
	})

function displayError(message) {
	const messageBox = document.getElementById('messageBox')
	messageBox.textContent = message
	messageBox.classList.remove('alert-success', 'alert-danger')

	messageBox.classList.add('alert', 'alert-danger')

	setTimeout(() => {
		messageBox.textContent = ''
		messageBox.classList.remove('alert-success', 'alert-danger')
	}, 4000)
}
