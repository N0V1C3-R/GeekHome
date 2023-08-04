document.querySelectorAll('input').forEach(input => {
    input.addEventListener('focus', () => {
        input.parentElement.classList.add('input-focused');
    });

    input.addEventListener('blur', () => {
        if (input.value === '') {
            input.parentElement.classList.remove('input-focused');
        }
    });
});

document.getElementById('verifyButton').addEventListener('click', function() {
    const email = document.getElementById("email").value;
    const verificationCode = document.getElementById('verificationCode').value;

    // 发送验证码验证请求
    fetch('/verify', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({"email": email, "code": verificationCode })
    })
        .then(function(response) {
            if (response.ok) {
                // 验证码验证成功，可以进行自动登录等操作
                console.log('验证码验证成功');
            } else {
                // 验证码验证失败，显示错误信息
                throw new Error('验证码验证失败');
            }
        })
        .catch(function(error) {
            // 处理验证码验证失败的情况
            console.error(error);
        });
});

document.getElementById('cancelButton').addEventListener('click', function() {
    document.getElementById('verificationModal').style.display = 'none';
});

function checkPasswordStrength() {
    const password = document.getElementById('password').value;
    const bar1 = document.getElementById('bar1');
    const bar2 = document.getElementById('bar2');
    const bar3 = document.getElementById('bar3');

    bar1.classList.remove('password-meter-bar--weak');
    bar2.classList.remove('password-meter-bar--medium');
    bar3.classList.remove('password-meter-bar--strong');

    if (password.length < 6) {
        bar1.style.backgroundColor = 'lightgray';
        bar2.style.backgroundColor = 'lightgray';
        bar3.style.backgroundColor = 'lightgray';
    } else {
        const types = {
            uppercase: /[A-Z]/,
            lowercase: /[a-z]/,
            number: /[0-9]/,
            special: /!@#\$%\^&*_+<>?,\.:/
        };

        let count = 0;

        if (password.length >= 6) {
            for (const type in types) {
                if (types[type].test(password)) {
                    count++;
                }
            }
        }

        if (count === 1) {
            bar1.style.backgroundColor = 'red';
            bar2.style.backgroundColor = 'lightgray';
            bar3.style.backgroundColor = 'lightgray';
        } else if (count === 2) {
            bar1.style.backgroundColor = 'yellow';
            bar2.style.backgroundColor = 'yellow';
            bar3.style.backgroundColor = 'lightgray';
        } else if (count >= 3) {
            bar1.style.backgroundColor = 'limegreen';
            bar2.style.backgroundColor = 'limegreen';
            bar3.style.backgroundColor = 'limegreen';
        }
    }
}

function validatePassword() {
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const passwordMismatchError = document.getElementById('passwordMismatchError');

    if (password !== confirmPassword && password.length > 6 && confirmPassword.length > 0) {
        passwordMismatchError.style.display = 'inline';
    } else {
        passwordMismatchError.style.display = 'none';
    }
}

function closeVerificationModal() {
    document.getElementById("verificationModal").style.display = "none";
}

function genXRequestId() {
    const timestamp = Date.now();
    const random = Math.floor(Math.random() * 1000);
    return `${timestamp}-${random}`
}