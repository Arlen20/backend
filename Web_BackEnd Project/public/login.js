document.addEventListener('DOMContentLoaded', () => {

    document.getElementById('loginForm').addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value.trim();

        if (username === '' || password === '') {
            displayMessage('Please fill in all fields.');
            return;
        }
        const formData = {
            username: username,
            password: password
        };
    
        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            if (response.ok) {
                window.location.href = '/';
            } else {
                const errorMessage = await response.text();
                console.error('Login failed:', errorMessage);
                displayError(errorMessage);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    });
});


function displayError(message) {
    const messageBox = document.getElementById('messageBox');
    messageBox.textContent = message;
    messageBox.classList.remove('alert-success', 'alert-danger');
    
    messageBox.classList.add('alert', 'alert-danger');

    setTimeout(() => {
        messageBox.textContent = '';
        messageBox.classList.remove('alert-success', 'alert-danger');
    }, 4000);
}