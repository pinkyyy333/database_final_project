<?php
session_start();

// 清除所有 session 變數
$_SESSION = array();

// 如果使用 cookie 儲存 session id，刪除 cookie
if (ini_get("session.use_cookies")) {
    $params = session_get_cookie_params();
    setcookie(session_name(), '', time() - 42000,
        $params["path"], $params["domain"],
        $params["secure"], $params["httponly"]
    );
}

// 最後銷毀 session
session_destroy();

// 回傳成功響應
http_response_code(200);
echo json_encode(['message' => 'Logout successful']);
exit;
