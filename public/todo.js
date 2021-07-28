(function($) {
  'use strict';
  $(function() {
    var todoListItem = $('.todo-list');
    var todoListInput = $('.todo-list-input');
    
    // Add
    $('.todo-list-add-btn').on("click", function(event) {
      event.preventDefault();
      
      var item = $(this).prevAll('.todo-list-input').val();
      if (item) {
        $.post("/todos", {name:item}, addItem)
        todoListInput.val("");
      }
    });

  
    let addItem = function(item) {
      if(item.completed) {
        todoListItem.append(`<li class="completed" id=`+item.id+`><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked'/>` + item.name + `<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>`);
      } else {
        todoListItem.append("<li id="+item.id+"><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
      }  
    }

    // 페이지 로딩시 서버로 부터 데이터를 읽어옴
    $.get("/todos", function(itemes) {
      itemes.forEach(element => {
        addItem(element)
      });
    })

    // Complete
    todoListItem.on('change', '.checkbox', function() {
      let id = $(this).closest('li').attr('id');      
      let $self = $(this);
      let complete = true;
      if($(this).attr('checked')) {
        complete = false;
      }
      $.get("complete-todo/"+id+"?complete"+complete, 
        function (data) {
          if (complete) {
            $self.attr('checked', 'checked');
          } else {  
            $self.removeAttr('checked');
          }
          $self.closest("li").toggleClass('completed');
        } 
      );
    });

    // Delete
    todoListItem.on('click', '.remove', function() {
      let id = $(this).closest('li').attr('id');      
      let $self = $(this);
      $.ajax({
        type: "DELETE",
        url: "todos/"+id,
        success: function (data) {
          if(data.success) {
            $self.parent().remove();
          }
        }
      });
    });
  });
})(jQuery);