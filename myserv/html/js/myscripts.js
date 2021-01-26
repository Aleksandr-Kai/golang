function onPageLoaded() {
    LoadAlbums();
    document.getElementById('myModal').addEventListener('show.bs.modal', function(event) {
        var xhr= new XMLHttpRequest();
        var button = event.relatedTarget
        var url = button.getAttribute('data-href')
        xhr.open('GET', url, true);
        xhr.onreadystatechange = function() {
            if (this.readyState!==4) return;
            if (this.status!==200) return;
            document.getElementById('mb').innerHTML= this.responseText;
            document.getElementById('list').focus();
            var title = $('h1#gallery-title').text();
            $('h1#modal-title').text(title);
            $('.gallery-buttons').hide();
            UpdateFancyBox();
        };
        xhr.send();
    })
    document.getElementById('myModal').addEventListener('hide.bs.modal', function(event) {
        document.getElementById('mb').innerHTML= "";
    })
    $('a#card-btn').on('click', function(event){
        $('.album-list').load($(this).attr('data-href'));
    })
    
}

function UpdateFancyBox(){
    $('[data-fancybox="gallery"]').fancybox({
        buttons: [
            "zoom",
            "share",
            "slideShow",
            "fullScreen",
            "download",
            "thumbs",
            "close"
          ],
        afterShow : function(instance, current) {
          var src =  current.src.replace('size=m', 'size=l');
      
          console.info('current.src: ' + current.src + '; download src: ' + src);
      
          $("[data-fancybox-download]").attr('href', src);
        },
        loop : true,
        animationEffect: "fade",
    });
}

function LoadAlbums(){
    $('.album-list').load('home?get_content=album-list', function(){
        $('img[src=""]').each(function(index){
            $(this).parent().parent().hide();
        });
        /*
        $('.card-img-top').on('click', function(){
            var src = $(this).attr('src').replace('size=s', 'size=m');
            $.fancybox.open('<div class="message"><img src="' + src + '></div>');
        });*/
    });
    
}

function GoBack(){
    LoadAlbums();
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}

function ScrollToTop(){
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}

function OpenGallery(gallery){
    $('.album-list').load(gallery, function(){UpdateFancyBox();});
    $([document.documentElement, document.body]).animate({
        scrollTop: 0
    }, 200);
}

if (document.readyState === 'complete' ||
    (document.readyState !== 'loading' && !document.documentElement.doScroll)) {
    onPageLoaded()
} else {
document.addEventListener('DOMContentLoaded', () => {
    onPageLoaded()
});
}

