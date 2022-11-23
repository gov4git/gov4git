---
name: Join the community
about: Use this form to join the community
title: "I'd like to join this project's community"
labels: ['community']
assignees: 'petar'
body:
  - type: markdown
    attributes:
      value: |
        Thank you for your interest in our community!
        Use the form below to join our community governance.
  - type: input
    id: contributor_public_url
    attributes:
      label: Your public repo
      description: The URL of your gov4git public repo
      placeholder: e.g. https://github.com/petar/gov4git.public.git
    validations:
      required: true
  - type: input
    id: contributor_public_branch
    attributes:
      label: Your public branch
      description: The branch within your gov4git public repo that holds your identity
      placeholder: e.g. main
    validations:
      required: true
  - type: input
    id: contributor_email
    attributes:
      label: Your email (optional)
      description: Share your email, if you want to receive community updates
      placeholder: e.g. petar@protocol.ai
    validations:
      required: false
