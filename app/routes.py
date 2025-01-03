from flask import Blueprint, request, jsonify
bp = Blueprint('api', __name__, url_prefix='/api')

@bp.route('/ping', methods=['GET'])
def get_ping():
    return jsonify(dict(data='pong'))

def register_routes(app):
    app.register_blueprint(bp)