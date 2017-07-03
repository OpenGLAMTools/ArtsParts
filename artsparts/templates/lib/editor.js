$(document).ready(function () {
    $('#image').cropper({

        viewMode: 1,
        dragMode: 'move',
        preview: '#preview',
        cropBoxMovable: false,
        guides: false,
        center: false,
        background: false,
        crop: function (e) {
            // Output the result data for cropping image.
            console.log(e.x);
            console.log(e.y);
            console.log(e.width);
            console.log(e.height);
            console.log(e.rotate);
            console.log(e.scaleX);
            console.log(e.scaleY);
        }
    });

});

var app = new Vue({
   delimiters: ['[[', ']]'],
    el: '#editor',
    data: {
        tweettext: "@OpenArtsParts"
    },
    computed: {
        charsRemain: function () {
            return 140 - this.tweettext.length
        },
        tooMuchChars: function (){
            return (140 < this.tweettext.length)
        }
    }
})