// axios({
//   method: 'post',
//   url: '/login',
//   data: {
//     firstName: 'Finn',
//     lastName: 'Williams',
//   },
// });

function DeleteBook(bookId) {
  axios({
    method: 'delete',
    url: `/api/books/${bookId}`,
    data: {
      firstName: 'Finn',
      lastName: 'Williams',
    },
  })
    .then(() => location.reload())
    .catch(() => console.log('error'));
}

function updatePage() {
  let id = document.location.pathname.split('/')[2];
  axios({
    method: 'get',
    url: `/api/books/${id}`,
  }).then((res) => {
    let data = res.data;
    this.document.querySelector('#Isbn').value = data.Isbn;
    this.document.querySelector('#Title').value = data.Title;
    this.document.querySelector('#Price').value = data.Price;
  });
  // ParseAuthor('Author-updateBook');
  let title = document.querySelector('#update-title');
  title.textContent;
}

function updateBook() {
  let id = document.location.pathname.split('/')[2];
  axios({
    method: 'put',
    url: `/api/books/${id}`,
    data: {
      Isbn: this.document.querySelector('#Isbn').value,
      Title: this.document.querySelector('#Title').value,
      Price: parseInt(this.document.querySelector('#Price').value),
    },
  })
    .then(() => location.replace('/home'))
    .catch(() => console.log('error'));
}

function ParseAuthor(selectID) {
  axios({
    method: 'get',
    url: `/api/authors`,
  }).then((res) => {
    let data = res.data;
    let sel = this.document.querySelector(`#${selectID}`);
    let inner = '<option>Select Author</option>';
    for (let i of data) {
      console.log('append');
      inner += `<option value="${i.ID}">${i.First}</option>`;
    }
    sel.innerHTML = inner;
  });
}

function addBook() {
  let author = parseInt(this.document.querySelector('#Author-newBook').value);
  if (!author) alert('please select an author');
  else {
    axios({
      method: 'post',
      url: `/api/books`,
      data: {
        Isbn: this.document.querySelector('#Isbn-newBook').value,
        Title: this.document.querySelector('#Title-newBook').value,
        Price: parseInt(this.document.querySelector('#Price-newBook').value),
        AuthorID: author,
      },
    })
      .then(() => location.reload())
      .catch(() => console.log('error'));
  }
}

function addAuthor() {
  axios({
    method: 'post',
    url: `/api/authors`,
    data: {
      First: this.document.querySelector('#First').value,
      Last: this.document.querySelector('#Last').value,
    },
  })
    .then(() => location.reload())
    .catch(() => console.log('error'));
}
