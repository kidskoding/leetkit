def is_strong_password(password):
    return (
        len(password) >= 8
        and any(ch.isupper() for ch in password)
        and any(ch.islower() for ch in password)
        and any(ch.isdigit() for ch in password)
        and all(not ch.isspace() for ch in password)
    )
