const form = document.getElementById("create_meme_form");
const fileInput = document.getElementById("image_file");

fileInput.addEventListener('change', () => {
  form.submit();
});

window.addEventListener('paste', e => {
  fileInput.files = e.clipboardData.files;
});