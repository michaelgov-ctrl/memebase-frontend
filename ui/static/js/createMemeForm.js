const form = document.getElementById("new_document_attachment");
const fileInput = document.getElementById("document_attachment_doc");

fileInput.addEventListener('change', () => {
  form.submit();
});

window.addEventListener('paste', e => {
  fileInput.files = e.clipboardData.files;
});