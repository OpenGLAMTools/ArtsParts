$(document).ready(function () {
    var app = new Vue({
        delimiters: ['[[', ']]'],
        el: '#editor',
        data: {
            tweettext: "@OpenArtsParts",
            apx: 0,
            apy: 0,
            apwidth: 0,
            apheight: 0
        },
        computed: {
            charsRemain: function () {
                return 140 - this.tweettext.length
            },
            tooMuchChars: function () {
                return (140 < this.tweettext.length)
            }
        },
        methods:{
            zoomIn: function(){
                 $('#image').cropper('zoom',0.25);
            },
            zoomOut: function(){
                 $('#image').cropper('zoom',-0.25);
            }
        }
    })
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
            cdata = $('#image').cropper('getCanvasData');
            w = cdata.naturalWidth;
            h = cdata.naturalHeight;
            app.apx = e.x / w;
            app.apy = e.y / h;
            app.apwidth = e.width / w;
            app.apheight = e.height / h;
            console.log("---------")
            console.log(e.x);
            console.log("--> " + e.x / cdata.naturalWidth);
            console.log(e.y);
            console.log("--> " + e.y / cdata.naturalHeight);
            console.log(e.width);
            console.log(e.height);
            console.log(e.rotate);
            console.log(e.scaleX);
            console.log(e.scaleY);
        }
    });
});

