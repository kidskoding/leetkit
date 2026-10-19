def email_domains(emails):
    return {email.rpartition("@")[2].lower() for email in emails}
