from flask import Flask

def create_app():
    app = Flask(__name__)
    # app.config.from_object('config.Config')

    with app.app_context():
        from .routes import register_routes
        register_routes(app)
        
    return app