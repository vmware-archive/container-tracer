# coding: utf-8

from __future__ import absolute_import

from flask import json
from six import BytesIO

from swagger_server.test import BaseTestCase


class TestDefaultController(BaseTestCase):
    """DefaultController integration test stubs"""

    def test_system_call_delete(self):
        """Test case for system_call_delete

        
        """
        response = self.client.open(
            '/v1/system_call',
            method='DELETE')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))

    def test_system_call_get(self):
        """Test case for system_call_get

        
        """
        response = self.client.open(
            '/v1/system_call',
            method='GET')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))

    def test_system_call_post(self):
        """Test case for system_call_post

        
        """
        response = self.client.open(
            '/v1/system_call',
            method='POST')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))

    def test_system_call_put(self):
        """Test case for system_call_put

        
        """
        response = self.client.open(
            '/v1/system_call',
            method='PUT')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    import unittest
    unittest.main()
