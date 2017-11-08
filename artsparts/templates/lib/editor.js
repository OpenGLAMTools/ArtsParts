$(document).ready(function () {
    var MAXCHARS = 280 
    var app = new Vue({
        delimiters: ['[[', ']]'],
        el: '#editor',
        data: {
            artpart: {
                tweettext: "http://artsparts.de"+PermanentLink+" #ArtsParts "+HashTag,
                x: 0,
                y: 0,
                width: 0,
                height: 0
            }
        },
        computed: {
            charsRemain: function () {
                return MAXCHARS - this.artpart.tweettext.length
            },
            tooMuchChars: function () {
                return (MAXCHARS < this.artpart.tweettext.length)
            }
        },
        methods: {
            zoomIn: function () {
                $('#image').cropper('zoom', 0.25);
            },
            zoomOut: function () {
                $('#image').cropper('zoom', -0.25);
            },
            createArtpart: function () {
                $('.ui.basic.modal').modal('show');
                this.$http.post('/artpart' + URIPath, this.artpart).then(response => {
                    //$('#artworkedit').modal('hide')

                    // success callback
                    console.log("Artwork is safed");
                    window.location.href = '/artwork' + URIPath;
                }, response => {
                    // error callback
                     console.log("There was an error");
                });
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
            app.artpart.x = e.x / w;
            app.artpart.y = e.y / h;
            app.artpart.width = e.width / w;
            app.artpart.height = e.height / h;
            
        }
    });
});

