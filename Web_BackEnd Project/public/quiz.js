$(document).ready(function () {
	let quizQuestions = []

	// Fetch quiz questions from the server
	$.ajax({
		url: '/quiz/questions',
		method: 'GET',
		success: function (response) {
			quizQuestions = response.quizQuestions
			renderQuizQuestions(quizQuestions)
		},
		error: function (error) {
			console.error('Error fetching quiz questions:', error)
		},
	})

	$('#btn').click(function () {
		let answers = []
		for (let i = 0; i < quizQuestions.length; i++) {
			const selectedAnswers = getSelectedAnswers('q' + i)
			console.log(`Question ${i + 1}: Selected answers: ${selectedAnswers}`)
			answers.push({
				questionId: quizQuestions[i]._id,
				answer: selectedAnswers,
			})
		}

		const userId = $('#userId').val()
		if (!userId) {
			alert('You need to be logged in to submit the quiz.')
			return
		}

		console.log('Submitting answers:', answers)

		$.ajax({
			url: '/quiz/submit',
			method: 'POST',
			contentType: 'application/json',
			data: JSON.stringify({
				userId,
				answers,
			}),
			success: function (response) {
				console.log('Success:', response)
				alert('Quiz results saved successfully.')
				document.getElementById('total').innerHTML = response.totalScore + '%'
			},
			error: function (error) {
				console.error('Error:', error)
				alert('Failed to save quiz results.')
			},
		})
	})

	function renderQuizQuestions(questions) {
		const quizForm = $('#quizForm')
		quizForm.empty() // Clear any existing questions to avoid duplication
		questions.forEach((question, index) => {
			const questionDiv = $('<div class="questionbar"></div>')
			const questionTitle = $('<h4></h4>').text('Question ' + (index + 1))
			const questionText = $('<p></p>').text(question.question)

			questionDiv.append(questionTitle, questionText)

			if (question.correctAnswer.length > 1) {
				question.options.forEach(option => {
					const label = $('<label></label>')
					const input = $(
						'<input class="form-check-input m-2" type="checkbox">'
					)
						.attr('name', 'q' + index)
						.val(option)
					label.append(input, option)
					questionDiv.append(label, '<br>')
				})
			} else {
				question.options.forEach(option => {
					const label = $('<label></label>')
					const input = $('<input class="form-check-input m-2" type="radio">')
						.attr('name', 'q' + index)
						.val(option)
					label.append(input, option)
					questionDiv.append(label, '<br>')
				})
			}

			quizForm.append(questionDiv)
		})
	}

	function getSelectedAnswers(name) {
		let checkboxes = document.getElementsByName(name)
		let selectedAnswers = []
		for (let i = 0; i < checkboxes.length; i++) {
			if (checkboxes[i].checked) {
				selectedAnswers.push(checkboxes[i].value)
			}
		}
		return selectedAnswers.join(', ')
	}
})
