var createDirForm = document.getElementById("create_dir_form");
var createDirDialogOpenButton = document.getElementById("create_dir_dialog_open_button");
var createDirDialogCloseButton = document.getElementById("create_dir_dialog_close_button");

createDirDialogOpenButton.onclick = function() {
    var createDirDialog=document.getElementById("create_dir_dialog");
    createDirDialog.open = true;
}

createDirDialogCloseButton.onclick = function() {
    var createDirDialog=document.getElementById("create_dir_dialog");
    createDirDialog.open = false;
}

createDirForm.onsubmit = async function(e) {
    e.preventDefault();

    var route = "/mkdir"
    var urlParameters = new URLSearchParams(window.location.search);
    var whichdirParam = urlParameters.get('whichdir');
    if (!whichdirParam) {
        whichdirParam = "";
    }
    var whichdir = "whichdir=" + whichdirParam;
    var createDirInput = document.getElementById("create_dir_input");
    var newdir = "newdir=" + createDirInput.value;
    var uri = route + "?" + whichdir + "&" + newdir;

    try {
        var response = await fetch(uri, {
            method: 'POST'
        });
        if (response.status == 409) {
            alert( createDirInput.value + " already created. Redirecting you there");
            return;
        } else if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
            return;
        } else {
            console.log(response);
        }
    } catch (err) {
        console.log(err);
        alert("is broken, let vvesley know");
        return
    }

    var newlocation = "";
    if (!whichdirParam) {
        newlocation = "/?whichdir="+createDirInput.value;
    } else {
        newlocation = "/?whichdir="+whichdirParam+"/"+createDirInput.value;
    }
    location.replace(newlocation);
}
